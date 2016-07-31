package main

import (
	"flag"
	"fmt"
	"os"
)

type programOptions struct {
	Region          string
	Username        string
	Password        string
	NetworkDomainID string
	RuleName        string
}

func (options *programOptions) Parse() bool {
	options.Username = os.Getenv("DD_COMPUTE_USER")
	if len(options.Username) == 0 {
		fmt.Println("The 'DD_COMPUTE_USER' environment variable is not defined. Set this to your CloudControl user name.")

		return false
	}
	options.Password = os.Getenv("DD_COMPUTE_PASSWORD")
	if len(options.Password) == 0 {
		fmt.Println("The 'DD_COMPUTE_PASSWORD' environment variable is not defined. Set this to your CloudControl password.")

		return false
	}

	flag.StringVar(&options.Region, "region", "AU", "The CloudControl region to use.")
	flag.StringVar(&options.NetworkDomainID, "network-domain", "", "The Id of the CloudControl network domain.")
	flag.StringVar(&options.RuleName, "rule-name", "", "The name of the firewall rule to configure.")
	flag.Parse()

	if len(options.Region) == 0 {
		fmt.Println("Must specify the Cloud Control region to use.")
		flag.Usage()

		return false
	}
	if len(options.NetworkDomainID) == 0 {
		fmt.Println("Must specify the Cloud Control network domain to use.")
		flag.Usage()

		return false
	}
	if len(options.RuleName) == 0 {
		fmt.Println("Must specify the name of the target firewall rule.")
		flag.Usage()

		return false
	}

	return true
}
