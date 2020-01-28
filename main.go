package main


import (
	"bitbucket.org/accezz-io/terraform-provider-symcsc/provider"
	"github.com/hashicorp/terraform/plugin"
)

var RateLimitSleepDuration = 5

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
