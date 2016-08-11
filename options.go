package main

import (
	"flag"
	"fmt"
	"os"
)

type programOptions struct {
	Region        string
	Username      string
	Password      string
	NetworkDomain string
	DataCenter    string
	RuleName      string
	Delete        bool
}

func (options *programOptions) Parse() bool {
	options.Username = os.Getenv("DD_COMPUTE_USER")
	if isEmpty(options.Username) {
		fmt.Println("The 'DD_COMPUTE_USER' environment variable is not defined. Set this to your CloudControl user name.")

		return false
	}
	options.Password = os.Getenv("DD_COMPUTE_PASSWORD")
	if isEmpty(options.Password) {
		fmt.Println("The 'DD_COMPUTE_PASSWORD' environment variable is not defined. Set this to your CloudControl password.")

		return false
	}

	flag.StringVar(&options.Region, "region", "AU", "The CloudControl region to use.")
	flag.StringVar(&options.NetworkDomain, "network", "", "The Id or name of the CloudControl network domain.")
	flag.StringVar(&options.DataCenter, "dc", "", "If the network domain's name is specified, the Id of the data center (e.g. AU9, NA1) in which the network domain is located.")
	flag.StringVar(&options.RuleName, "rule", "", "The name of the firewall rule to configure.")
	flag.BoolVar(&options.Delete, "delete", false, "Delete the firewall rule, if it exists.")
	flag.Parse()

	if isEmpty(options.Region) {
		fmt.Println("Must specify the Cloud Control region to use.")
		flag.Usage()

		return false
	}
	if isEmpty(options.NetworkDomain) {
		fmt.Println("Must specify the Cloud Control network domain to use.")
		flag.Usage()

		return false
	} else if !isUUID(options.NetworkDomain) && isEmpty(options.DataCenter) {
		fmt.Println("If the network domain is specified by name (instead of Id), must specify the Cloud Control data centre that contains it.")
		flag.Usage()

		return false
	}
	if isEmpty(options.RuleName) {
		fmt.Println("Must specify the name of the target firewall rule.")
		flag.Usage()

		return false
	}

	return true
}
