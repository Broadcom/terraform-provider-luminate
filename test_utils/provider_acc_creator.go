package test_utils

import (
	"context"
	luminateFrameworkProvider "github.com/Broadcom/terraform-provider-luminate/framework_provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
)

func CreateProtoV6ProviderFactories(sdkTestProvider *schema.Provider) map[string]func() (tfprotov6.ProviderServer, error) {
	ctx := context.Background()
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(
		ctx,
		sdkTestProvider.GRPCProvider,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
		providerserver.NewProtocol6(luminateFrameworkProvider.NewLuminateFrameworkProvider(sdkTestProvider)),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return map[string]func() (tfprotov6.ProviderServer, error){
		"luminate": func() (tfprotov6.ProviderServer, error) {
			return muxServer.ProviderServer(), nil
		},
	}
}
