// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

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
			ValidateFunc: utils.ValidateBool,
			Default:      true,
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
			Type:        schema.TypeString,
			Description: "The application effective DNS address that exposes the application",
			Computed:    true,
		},
		"luminate_address": {
			Type:        schema.TypeString,
			Description: "Application DNS address using Symantec ZTNA subdomain (i.e., testapp.acme.luminatesec.com)",
			Computed:    true,
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

// The management server appends the application type's default port to an internal address
// submitted without one, so a configured "tcp://1.1.1.1" is returned as "tcp://1.1.1.1:22".
// Removing it again keeps state equal to the configured value whether or not the server
// applies that defaulting.
func stripDefaultPort(address string, defaultPort string) string {
	if _, port := utils.ExtractIPAndPort(address); port != defaultPort {
		return address
	}

	return strings.TrimSuffix(address, ":"+defaultPort)
}

// suppressDefaultPortDiff equates an address carrying the type's default port with the same
// address omitting it, so neither form is reported as a change.
func suppressDefaultPortDiff(defaultPort string) schema.SchemaDiffSuppressFunc {
	return func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		if oldValue == "" {
			return false
		}
		if oldValue == newValue {
			return true
		}

		oldAddress, oldPort := utils.ExtractIPAndPort(oldValue)
		newAddress, newPort := utils.ExtractIPAndPort(newValue)
		if oldAddress != newAddress {
			return false
		}
		if (oldPort == "" && newPort == defaultPort) || (oldPort == defaultPort && newPort == "") {
			return true
		}

		return oldPort == newPort
	}
}
