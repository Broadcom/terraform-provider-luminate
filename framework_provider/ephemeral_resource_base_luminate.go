package framework_provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"

	"github.com/Broadcom/terraform-provider-luminate/service"
)

type BaseLuminateEphemeralResource struct {
	client *service.LuminateService
}

func (r *BaseLuminateEphemeralResource) Configure(_ context.Context, request ephemeral.ConfigureRequest, response *ephemeral.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*service.LuminateService)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *service.LuminateService, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	r.client = client
}
