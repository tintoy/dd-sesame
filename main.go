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

	// If user specified network domain by name instead of Id, resolve it now.
	err = findTargetNetworkDomainID(apiClient, &options)
	if err != nil {
		panic(err)
	}

	var targetRule *compute.FirewallRule
	targetRule, err = findFirewallRuleByName(apiClient, options)
	if err != nil {
		panic(err)
	}
	if targetRule != nil {
		if !options.Delete && targetRule.Source.IPAddress != nil && targetRule.Source.IPAddress.Address == externalIP {
			fmt.Printf("Firewall rule '%s' ('%s') is already up-to-date.\n", options.RuleName, targetRule.ID)

			return
		}

		err := apiClient.DeleteFirewallRule(targetRule.ID)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Deleted existing firewall rule '%s' (Id = '%s').\n", options.RuleName, targetRule.ID)
	}

	if options.Delete {
		return
	}

	ruleConfiguration := compute.FirewallRuleConfiguration{
		Name:            options.RuleName,
		NetworkDomainID: options.NetworkDomain,
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

func findTargetNetworkDomainID(apiClient *compute.Client, options *programOptions) (err error) {
	if isUUID(options.NetworkDomain) {
		return
	}

	fmt.Printf("Looking up network domain '%s' in data centre '%s'...\n", options.NetworkDomain, options.DataCenter)

	var networkDomain *compute.NetworkDomain
	networkDomain, err = apiClient.GetNetworkDomainByName(options.NetworkDomain, options.DataCenter)
	if err != nil {
		return
	}
	if networkDomain == nil {
		err = fmt.Errorf("No network domain named '%s' was found in data centre '%s'", options.NetworkDomain, options.DataCenter)

		return
	}

	fmt.Printf("Found network domain '%s' (Id = '%s').\n", networkDomain.Name, networkDomain.ID)

	options.NetworkDomain = networkDomain.ID

	return
}

func findFirewallRuleByName(apiClient *compute.Client, options programOptions) (firewallRule *compute.FirewallRule, err error) {
	page := compute.DefaultPaging()
	page.PageSize = 50
	for firewallRule == nil {
		var firewallRules *compute.FirewallRules
		firewallRules, err = apiClient.ListFirewallRules(options.NetworkDomain, page)
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
