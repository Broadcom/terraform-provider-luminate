// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/Broadcom/terraform-provider-luminate/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
