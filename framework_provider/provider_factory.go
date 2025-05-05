// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package framework_provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateProviderServer(sdkProvider *schema.Provider) (func() tfprotov6.ProviderServer, error) {
	ctx := context.Background()
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(
		ctx,
		sdkProvider.GRPCProvider,
	)
	if err != nil {
		return nil, err
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
		providerserver.NewProtocol6(NewLuminateFrameworkProvider(sdkProvider)),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil

}
