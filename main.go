package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	luminateFrameworkProvider "github.com/Broadcom/terraform-provider-luminate/framework_provider"
	luminateSdkProvider "github.com/Broadcom/terraform-provider-luminate/provider"
)

func main() {
	withDebug := flag.Bool("debug", false, "runs the provider with debug to allows debuggers such as delve to attach to it")
	flag.Parse()

	var opts []tf6server.ServeOpt
	if withDebug != nil && *withDebug {
		opts = append(opts, tf6server.WithManagedDebug())
	}

	ctx := context.Background()
	sdkProvider := luminateSdkProvider.Provider()
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(
		context.Background(),
		sdkProvider.GRPCProvider,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
		providerserver.NewProtocol6(luminateFrameworkProvider.NewLuminateFrameworkProvider(sdkProvider)),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = tf6server.Serve("broadcom.com/broadcom/luminate", muxServer.ProviderServer, opts...)

	if err != nil {
		log.Fatalln(err.Error())
	}
}
