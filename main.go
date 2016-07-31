package main

import (
	"compute-api/compute"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var options = programOptions{}

func main() {
	options := programOptions{}
	if !options.Parse() {
		os.Exit(2)
	}

	fmt.Println("Detecting external IP address...")

	externalIP, err := getExternalIP()
	if err != nil {
		panic(err)
	}

	fmt.Printf("External IP address is '%s'.\n", externalIP)

	apiClient := compute.NewClient(options.Region, options.Username, options.Password)

	var targetRule *compute.FirewallRule
	targetRule, err = findFirewallRuleByName(apiClient, options)
	if err != nil {
		panic(err)
	}
	if targetRule != nil {
		if targetRule.Source.IPAddress != nil && targetRule.Source.IPAddress.Address == externalIP {
			fmt.Printf("Firewall rule '%s' ('%s') is already up-to-date.\n", options.RuleName, targetRule.ID)

			return
		}

		err := apiClient.DeleteFirewallRule(targetRule.ID)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Deleted existing firewall rule '%s' (Id = '%s').\n", options.RuleName, targetRule.ID)
	}

	ruleConfiguration := compute.FirewallRuleConfiguration{
		Name:            options.RuleName,
		NetworkDomainID: options.NetworkDomainID,
		Action:          "ACCEPT_DECISIVELY",
		Enabled:         true,
		Protocol:        "IP",
		IPVersion:       "IPv4",
	}
	ruleConfiguration.MatchSourceAddressAndPort(externalIP, nil)
	ruleConfiguration.MatchAnyDestination()
	ruleConfiguration.PlaceFirst()

	fmt.Printf("Creating firewall rule '%s'.\n", options.RuleName)
	firewallRuleID, err := apiClient.CreateFirewallRule(ruleConfiguration)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created firewall rule '%s' (Id = '%s').\n", options.RuleName, firewallRuleID)
}

func findFirewallRuleByName(apiClient *compute.Client, options programOptions) (firewallRule *compute.FirewallRule, err error) {
	page := compute.DefaultPaging()
	page.PageSize = 50
	for firewallRule == nil {
		var firewallRules *compute.FirewallRules
		firewallRules, err = apiClient.ListFirewallRules(options.NetworkDomainID, page)
		if err != nil {
			return
		}
		if firewallRules.IsEmpty() {
			return
		}

		for _, rule := range firewallRules.Rules {
			if rule.Name == options.RuleName {
				firewallRule = &rule

				return
			}
		}

		page.Next()
	}

	return
}

func getExternalIP() (externalIP string, err error) {
	var response *http.Response
	response, err = http.Get("http://ifconfig.co/")
	if err != nil {
		return
	}

	defer response.Body.Close()

	var responseBody []byte
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	externalIP = strings.TrimSpace(
		string(responseBody),
	)

	return
}
