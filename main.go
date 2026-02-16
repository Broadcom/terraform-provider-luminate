// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"log"

	luminateFrameworkProvider "github.com/Broadcom/terraform-provider-luminate/framework_provider"
	luminateSdkProvider "github.com/Broadcom/terraform-provider-luminate/provider"
)

func main() {
	withDebug := flag.Bool("debug", true, "runs the provider with debug to allows debuggers such as delve to attach to it")
	flag.Parse()

	var opts []tf6server.ServeOpt
	if withDebug != nil && *withDebug {
		opts = append(opts, tf6server.WithManagedDebug())
	}

	sdkProvider := luminateSdkProvider.Provider()
	providerServer, err := luminateFrameworkProvider.CreateProviderServer(sdkProvider)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = tf6server.Serve("broadcom.com/broadcom/luminate", providerServer, opts...)

	if err != nil {
		log.Fatalln(err.Error())
	}
}
