package provider

import (
	"errors"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/validation"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform/helper/schema"
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
		Schema: dnsSchema,
		Create: resourceCreateDNSServer,
		Read:   resourceReadDNSServer,
		Update: resourceUpdateDNSServer,
		Delete: resourceDeleteDNSServer,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateDNSServer(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] LUMINATE CREATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	newApp := extractDNSServerFields(d)

	app, err := client.Applications.CreateApplication(newApp)
	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return err
	}

	d.SetId(app.ID)
	setDNSServerFields(d, app)

	return resourceReadDNSServer(d, m)
}

func resourceReadDNSServer(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] LUMINATE READ APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	app, err := client.Applications.GetApplicationById(d.Id())

	if err != nil {
		if err.Error() == "403 Forbidden" {
			d.SetId("")
			return nil
		}
		return err
	}

	if app == nil {
		d.SetId("")
		return nil
	}

	d.SetId(app.ID)

	app.SiteID = d.Get("site_id").(string)
	setDNSServerFields(d, app)

	return nil
}

func resourceUpdateDNSServer(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE UPDATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	app := extractDNSServerFields(d)

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(app, app.SiteID)
	if err != nil {
		return err
	}

	updApp.SiteID = app.SiteID
	setDNSServerFields(d, updApp)

	return resourceReadDNSServer(d, m)
}

func resourceDeleteDNSServer(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE DELETE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	err := client.Applications.DeleteApplication(d.Id())
	if err != nil {
		return err
	}

	return resourceReadDNSServer(d, m)
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
