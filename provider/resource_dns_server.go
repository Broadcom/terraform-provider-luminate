package provider

import (
	"context"
	"errors"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LuminateDNSserver() *schema.Resource {
	dnsSchema := CommonApplicationSchema()

	dnsSchema["internal_address"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "Internal address of the application, accessible by connector",
	}

	dnsSchema["dns_settings"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"domain_suffixes": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.NoZeroValues,
					},
				},
			},
		},
	}

	return &schema.Resource{
		Schema:        dnsSchema,
		CreateContext: resourceCreateDNSServer,
		ReadContext:   resourceReadDNSServer,
		UpdateContext: resourceUpdateDNSServer,
		DeleteContext: resourceDeleteDNSServer,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateDNSServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] LUMINATE CREATE APP")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	newApp := extractDNSServerFields(d)

	app, err := client.Applications.CreateApplication(newApp)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.ID)
	setDNSServerFields(d, app)
	resourceReadDNSServer(ctx, d, m)
	return diagnostics
}

func resourceReadDNSServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] LUMINATE READ APP")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		err := errors.New("unable to cast Luminate service")
		return diag.FromErr(err)
	}

	app, err := client.Applications.GetApplicationById(d.Id())

	if err != nil {
		if err.Error() == "403 Forbidden" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if app == nil {
		d.SetId("")
		return nil
	}

	d.SetId(app.ID)

	app.SiteID = d.Get("site_id").(string)
	setDNSServerFields(d, app)

	return diagnostics
}

func resourceUpdateDNSServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE UPDATE APP")
	var diagnostics diag.Diagnostics

	client, ok := m.(*service.LuminateService)
	if !ok {
		err := errors.New("unable to cast Luminate service")
		return diag.FromErr(err)
	}

	app := extractDNSServerFields(d)

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Applications.BindApplicationToSite(app, app.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	updApp.SiteID = app.SiteID
	setDNSServerFields(d, updApp)

	return diagnostics
}

func resourceDeleteDNSServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE DELETE APP")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		err := errors.New("unable to cast Luminate service")
		return diag.FromErr(err)
	}

	err := client.Applications.DeleteApplication(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diagnostics
}

func setDNSServerFields(d *schema.ResourceData, application *dto.Application) {
	d.Set("name", application.Name)
	d.Set("icon", application.Icon)
	d.Set("type", application.Type)
	d.Set("visible", application.Visible)
	d.Set("notification_enabled", application.NotificationsEnabled)
	d.Set("internal_address", application.InternalAddress)
	d.Set("external_address", application.ExternalAddress)
	if application.DnsSettings != nil {
		d.Set("dnsSettings", flattenDNSSettings(application.DnsSettings))
	}
}

func extractDNSServerFields(d *schema.ResourceData) *dto.Application {
	dnsSettings := extractDNSSettings(d)
	return &dto.Application{
		Name:                 d.Get("name").(string),
		Icon:                 d.Get("icon").(string),
		SiteID:               d.Get("site_id").(string),
		Type:                 "dns",
		Visible:              d.Get("visible").(bool),
		NotificationsEnabled: d.Get("notification_enabled").(bool),
		Subdomain:            d.Get("subdomain").(string),
		InternalAddress:      d.Get("internal_address").(string),
		ExternalAddress:      d.Get("external_address").(string),
		DnsSettings:          dnsSettings,
	}
}

func flattenDNSSettings(settings *dto.DnsSettings) []interface{} {
	if settings == nil {
		return []interface{}{}
	}

	k := map[string]interface{}{
		"domainSuffixes": settings.DomainSuffixes,
	}

	return []interface{}{k}
}

func extractDNSSettings(d *schema.ResourceData) *dto.DnsSettings {
	var dnsSettings *dto.DnsSettings
	if v, ok := d.GetOk("dns_settings"); ok {
		for _, element := range v.([]interface{}) {
			elem := element.(map[string]interface{})

			var domainSuffixes []string

			if domainSuffix, ok := elem["domain_suffixes"].([]interface{}); ok {
				for _, suffix := range domainSuffix {
					domainSuffixes = append(domainSuffixes, suffix.(string))
				}
			}
			dnsSettings = &dto.DnsSettings{
				DomainSuffixes: domainSuffixes,
			}
		}
	}
	return dnsSettings
}
