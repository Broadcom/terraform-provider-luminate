// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const TCPApplicationInternalAddress = "tcp://luminate-dummy:22"

func LuminateTCPApplication() *schema.Resource {
	tcpSchema := CommonApplicationSchema()

	tcpSchema["target"] = &schema.Schema{
		Type:     schema.TypeList,
		MinItems: 1,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"address": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Application target address.",
					ValidateFunc: utils.ValidateString,
				},
				"ports": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
				"port_mapping": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
			},
		},
	}

	return &schema.Resource{
		Schema:        tcpSchema,
		CreateContext: resourceCreateTCPApplication,
		ReadContext:   resourceReadTCPApplication,
		UpdateContext: resourceUpdateTCPApplication,
		DeleteContext: resourceDeleteTCPApplication,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateTCPApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE CREATE APP")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}
	newApp := extractTCPApplicationFields(d)

	app, err := client.Applications.CreateApplication(newApp)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.ID)
	setTCPApplicationFields(d, app)

	resourceReadTCPApplication(ctx, d, m)
	return diagnostics
}

func resourceReadTCPApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] LUMINATE READ APP")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	app, err := client.Applications.GetApplicationById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if app == nil {
		d.SetId("")
		return nil
	}

	d.SetId(app.ID)
	app.SiteID = d.Get("site_id").(string)
	setTCPApplicationFields(d, app)
	return diagnostics
}

func resourceUpdateTCPApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE UPDATE APP")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	app := extractTCPApplicationFields(d)

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Applications.BindApplicationToSite(app, app.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	setTCPApplicationFields(d, updApp)

	resourceReadTCPApplication(ctx, d, m)
	return diagnostics
}

func resourceDeleteTCPApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE DELETE APP")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	err := client.Applications.DeleteApplication(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resourceReadTCPApplication(ctx, d, m)
	return diagnostics
}

func extractTCPApplicationFields(d *schema.ResourceData) *dto.Application {
	app := dto.Application{
		ID:                   d.Id(),
		Name:                 d.Get("name").(string),
		Icon:                 d.Get("icon").(string),
		SiteID:               d.Get("site_id").(string),
		CollectionID:         d.Get("collection_id").(string),
		Type:                 "tcp",
		Visible:              d.Get("visible").(bool),
		NotificationsEnabled: d.Get("notification_enabled").(bool),
		InternalAddress:      TCPApplicationInternalAddress,
		ExternalAddress:      d.Get("external_address").(string),
		Subdomain:            d.Get("subdomain").(string),
	}

	app.Targets = extractTCPTargets(d)

	return &app
}

func extractTCPTargets(d *schema.ResourceData) []dto.TCPTarget {
	var targets []dto.TCPTarget

	configTargets, ok := d.Get("target").([]interface{})
	if !ok {
		return targets
	}

	for _, v := range configTargets {
		log.Printf("[DEBUG] FOUND TCP Targets %v", v)
		vt := v.(map[string]interface{})

		p := vt["ports"].([]interface{})
		var ports []int32

		for _, port := range p {
			ports = append(ports, int32(port.(int)))
		}

		pm := vt["port_mapping"].([]interface{})
		var portMapping []int32

		for _, port := range pm {
			portMapping = append(portMapping, int32(port.(int)))
		}

		target := dto.TCPTarget{
			Address:     vt["address"].(string),
			Ports:       ports,
			PortMapping: portMapping,
		}
		targets = append(targets, target)
	}
	return targets
}

func setTCPApplicationFields(d *schema.ResourceData, application *dto.Application) {
	SetBaseApplicationFields(d, application)
	d.Set("collection_id", application.CollectionID)
	d.Set("luminate_address", application.LuminateAddress)
	log.Printf("[DEBUG] Settings TCP Targets")

	d.Set("target", flattenTCPTargets(application.Targets))
}

func flattenTCPTargets(targets []dto.TCPTarget) []interface{} {

	var flat []interface{}

	for _, c := range targets {
		t := map[string]interface{}{
			"address":      c.Address,
			"ports":        c.Ports,
			"port_mapping": c.PortMapping,
		}
		flat = append(flat, t)
	}
	return flat
}
