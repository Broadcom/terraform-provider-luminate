// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CommonApplicationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"site_id": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Site ID to which the application will be bound",
			ValidateFunc: utils.ValidateUuid,
		},
		"collection_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Collection ID to which the application will be assigned",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Name of the application",
			ValidateFunc: utils.ValidateApplicationName,
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "app type",
			ValidateFunc: utils.ValidateString,
			Computed:     true,
		},
		"icon": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Base64 representation of 128x128 icon",
			ValidateFunc: utils.ValidateString,
		},
		"visible": {
			Type:         schema.TypeBool,
			Optional:     true,
			Default:      true,
			ValidateFunc: utils.ValidateBool,
			Description:  "Indicates whether to show this application in the applications portal.",
		},
		"notification_enabled": {
			Type:         schema.TypeBool,
			Optional:     true,
			Default:      true,
			ValidateFunc: utils.ValidateBool,
			Description:  "Indicates whether notifications are enabled for this application.",
		},
		"subdomain": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The application DNS subdomain.",
		},
		"external_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"luminate_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
func SetBaseApplicationFields(d *schema.ResourceData, application *dto.Application) {
	d.Set("name", application.Name)
	if application.Icon != "" {
		d.Set("icon", application.Icon)
	}
	d.Set("visible", application.Visible)
	d.Set("notification_enabled", application.NotificationsEnabled)
	d.Set("external_address", application.ExternalAddress)
	d.Set("subdomain", application.Subdomain)
	d.Set("type", application.Type)
}
