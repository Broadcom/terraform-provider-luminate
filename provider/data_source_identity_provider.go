// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LuminateDataSourceIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"identity_provider_name": {
				Type:         schema.TypeString,
				Description:  "The identity provider name as configured in Luminate portal, if not specified local idp will be taken",
				Required:     true,
				ValidateFunc: utils.ValidateString,
			},
			"identity_provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: resourceReadIdentityProvider,
	}
}

func resourceReadIdentityProvider(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	identityProviderName := d.Get("identity_provider_name").(string)

	identityProviderId, err := client.IdentityProviders.GetIdentityProviderId(identityProviderName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(identityProviderName)
	d.Set("identity_provider_id", identityProviderId)

	return nil
}
