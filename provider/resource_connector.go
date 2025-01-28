// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LuminateConnector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "A descriptive name of the Connector",
				Required:     true,
				ValidateFunc: utils.ValidateString,
				ForceNew:     true,
			},
			"site_id": {
				Type:         schema.TypeString,
				Description:  "Site to bind the connector to",
				Required:     true,
				ValidateFunc: utils.ValidateUuid,
				ForceNew:     true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "The deployment type of the host running the Symantec ZTNA connector",
				Required:     true,
				ValidateFunc: validateConnectorType,
				ForceNew:     true,
			},
			"command": {
				Type:        schema.TypeString,
				Description: "Command for deploying Symantec ZTNA connector",
				Computed:    true,
			},
			"otp": {
				Type:        schema.TypeString,
				Description: "One time password for running Symantec ZTNA connector",
				Computed:    true,
			},
		},
		CreateContext: resourceCreateConnector,
		ReadContext:   resourceReadConnector,
		DeleteContext: resourceDeleteConnector,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateConnector(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}
	conOpts := extractConnectorFields(d)

	newCon, err := client.Connectors.CreateConnector(conOpts, conOpts.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	setConnectorFields(d, newCon)

	command, err := client.Connectors.GetConnectorCommand(newCon.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("command", command)

	return resourceReadConnector(ctx, d, m)
}

func resourceReadConnector(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*service.LuminateService)
	var diagnostics diag.Diagnostics
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	connector, err := client.Connectors.GetConnectorByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if connector == nil {
		d.SetId("")
		return nil
	}

	connector.SiteID = d.Get("site_id").(string)

	setConnectorFields(d, connector)

	return diagnostics
}

func resourceDeleteConnector(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	err := client.Connectors.DeleteConnector(d.Id())
	if err != nil {
		return diag.FromErr(errors.New("unable to delete connector"))
	}
	d.SetId("")

	return resourceReadConnector(ctx, d, m)
}

func validateConnectorType(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	cType, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type to be string"))
		return warns, errs
	}

	validTypes := []string{
		"linux",
		"kubernetes",
		"windows",
		"docker-compose",
	}

	if !utils.StringInSlice(validTypes, cType) {
		errs = append(errs, fmt.Errorf("connector type must be one of %v", validTypes))
	}
	return warns, errs
}

func setConnectorFields(d *schema.ResourceData, connector *dto.Connector) {
	d.SetId(connector.ID)
	d.Set("name", connector.Name)
	d.Set("type", connector.Type)
	d.Set("site_id", connector.SiteID)
	d.Set("otp", connector.OTP)
}

func extractConnectorFields(d *schema.ResourceData) *dto.Connector {
	return &dto.Connector{
		Name:   d.Get("name").(string),
		Type:   d.Get("type").(string),
		SiteID: d.Get("site_id").(string),
	}
}
