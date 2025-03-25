package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"

	luminateFrameworkProvider "github.com/Broadcom/terraform-provider-luminate/framework_provider"
	luminateSdkProvider "github.com/Broadcom/terraform-provider-luminate/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	withDebug := flag.Bool("debug", false, "runs the provider with debug to allows debuggers such as delve to attach to it")
	flag.Parse()

	var opts []tf5server.ServeOpt
	if withDebug != nil && *withDebug {
		opts = append(opts, tf5server.WithManagedDebug())
	}

	ctx := context.Background()
	sdkProvider := luminateSdkProvider.Provider()
	providers := []func() tfprotov5.ProviderServer{
		sdkProvider.GRPCProvider,
		providerserver.NewProtocol5(luminateFrameworkProvider.NewLuminateFrameworkProvider(sdkProvider)),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = tf5server.Serve("broadcom.com/broadcom/luminate", muxServer.ProviderServer, opts...)

	if err != nil {
		log.Fatalln(err.Error())
	}
}
