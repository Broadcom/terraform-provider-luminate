package main

import (
	"github.com/Broadcom/terraform-provider-luminate/provider"
	"github.com/hashicorp/terraform/plugin"
)

var RateLimitSleepDuration = 5

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
