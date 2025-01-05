// Copyright (c) Symantec ZTNA
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

func LuminateDataSourceAwsIntegration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"integration_name": {
				Type:         schema.TypeString,
				Description:  "The aws integration name as configured in Luminate portal",
				Required:     true,
				ValidateFunc: utils.ValidateString,
			},
			"integration_id": {
				Type:        schema.TypeString,
				Description: "The aws integration id",
				Computed:    true,
			},
		},
		ReadContext: resourceReadAwsIntegration,
	}
}

func resourceReadAwsIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client, ok := m.(*service.LuminateService)
	var diagnostics diag.Diagnostics
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	integrationName := d.Get("integration_name").(string)

	integrationId, err := client.IntegrationAPI.GetIntegrationId(integrationName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(integrationId)

	return diagnostics
}
