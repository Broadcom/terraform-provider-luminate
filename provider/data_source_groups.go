// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func LuminateDataSourceGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"identity_provider_id": {
				Type:         schema.TypeString,
				Description:  "The identity provider id",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"groups": {
				Type:        schema.TypeList,
				Description: "list of groups to include as part of this policy",
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},
			"group_ids": {
				Type:        schema.TypeList,
				Description: "A list containing the ids of the requested groups",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},

				Computed: true,
			},
		},
		ReadContext: resourceReadGroups,
	}
}

func resourceReadGroups(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diagnostic diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	var groupIds []string

	identityProviderId := d.Get("identity_provider_id").(string)
	groupNames := d.Get("groups").([]interface{})

	for _, groupName := range groupNames {
		groupId, err := client.Groups.GetGroupId(identityProviderId, groupName.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		groupIds = append(groupIds, groupId)
	}

	d.SetId(identityProviderId)
	d.Set("group_ids", groupIds)

	return diagnostic
}
