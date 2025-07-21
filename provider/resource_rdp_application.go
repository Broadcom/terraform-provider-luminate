// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
	"log"
	"regexp"
)

func LuminateRDPApplication() *schema.Resource {
	rdpSchema := CommonApplicationSchema()

	rdpSchema["internal_address"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ValidateFunc:     utils.ValidateString,
		Description:      "Internal address of the application, accessible by connector",
		DiffSuppressFunc: suppressExternalAddressUpdate,
	}

	rdpSchema["sub_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      string(sdk.SINGLE_MACHINE_ApplicationSubType),
		ValidateFunc: validateSubType,
		Description:  "rdp application sub type",
	}

	return &schema.Resource{
		Schema:        rdpSchema,
		CreateContext: resourceCreateRDPApplication,
		ReadContext:   resourceReadRDPApplication,
		UpdateContext: resourceUpdateRDPApplication,
		DeleteContext: resourceDeleteRDPApplication,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateRDPApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] LUMINATE CREATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	newApp := extractRDPApplicationFields(d)
	app, err := client.Applications.CreateApplication(newApp)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.ID)
	setRDPApplicationFields(d, app)

	return resourceReadRDPApplication(ctx, d, m)
}

func resourceReadRDPApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] LUMINATE READ APP")

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
	setRDPApplicationFields(d, app)

	return nil
}

func resourceUpdateRDPApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE UPDATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	app := extractRDPApplicationFields(d)

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Applications.BindApplicationToSite(app, app.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	updApp.SiteID = app.SiteID
	setRDPApplicationFields(d, updApp)

	return resourceReadRDPApplication(ctx, d, m)
}

func resourceDeleteRDPApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE DELETE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	err := client.Applications.DeleteApplication(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceReadRDPApplication(ctx, d, m)
}

func setRDPApplicationFields(d *schema.ResourceData, application *dto.Application) {
	SetBaseApplicationFields(d, application)
	d.Set("collection_id", application.CollectionID)
	d.Set("sub_type", application.SubType)
	d.Set("site_id", application.SiteID)
	d.Set("internal_address", application.InternalAddress)
	d.Set("luminate_address", application.LuminateAddress)
}

func extractRDPApplicationFields(d *schema.ResourceData) *dto.Application {

	// adding port to rdp app when provided without one
	internalAddress := d.Get("internal_address").(string)
	pattern := `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(internalAddress) {
		internalAddress = fmt.Sprintf("%s:%s", internalAddress, utils.DefaultRDPPort)
	}

	return &dto.Application{
		ID:                   d.Id(),
		Name:                 d.Get("name").(string),
		CollectionID:         d.Get("collection_id").(string),
		SubType:              d.Get("sub_type").(string),
		Icon:                 d.Get("icon").(string),
		SiteID:               d.Get("site_id").(string),
		Type:                 "rdp",
		Visible:              d.Get("visible").(bool),
		NotificationsEnabled: d.Get("notification_enabled").(bool),
		InternalAddress:      internalAddress,
		ExternalAddress:      d.Get("external_address").(string),
		Subdomain:            d.Get("subdomain").(string),
	}
}

// suppressExternalAddressUpdate will determine if needed another action (CRUD) from terraform, in will run after terraform plan is running
// if it returns false terraform will run another action when state != require value
func suppressExternalAddressUpdate(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue == "" {
		return false
	}
	if oldValue == newValue {
		return true
	}
	newAddress, newPort := utils.ExtractIPAndPort(newValue)
	oldAddress, oldPort := utils.ExtractIPAndPort(oldValue)
	if ((oldPort == "" && newPort == utils.DefaultRDPPort) || (oldPort == utils.DefaultRDPPort && newPort == "")) && (newAddress == oldAddress) {
		return true
	}
	if (newAddress != oldAddress) || (newPort != oldPort) {
		return false
	}
	return true
}
