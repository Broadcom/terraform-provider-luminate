package provider

import (
	"errors"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
)

func LuminateSegmentApplication() *schema.Resource {
	segmentAppSchema := CommonApplicationSchema()

	segmentAppSchema["segment_settings"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"original_ip": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Target ip",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: utils.ValidateString,
					},
				},
			},
		},
	}

	return &schema.Resource{
		Schema: segmentAppSchema,
		Create: resourceCreateSegmentApplication,
		Read:   resourceReadSegmentApplication,
		Update: resourceUpdateSegmentApplication,
		Delete: resourceDeleteSegmentApplication,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateSegmentApplication(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] LUMINATE CREATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	newApp := extractSegmentApplication(d)

	app, err := client.Applications.CreateApplication(newApp)

	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return err
	}

	d.SetId(app.ID)
	setSegmentApplicationFields(d, app, client.TenantBaseDomain)

	return resourceReadSegmentApplication(d, m)
}

func resourceReadSegmentApplication(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] LUMINATE READ APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	app, err := client.Applications.GetApplicationById(d.Id())
	if err != nil {
		return err
	}

	if app == nil {
		d.SetId("")
		return nil
	}

	app.SiteID = d.Get("site_id").(string)
	setSegmentApplicationFields(d, app, client.TenantBaseDomain)

	return nil
}

func resourceUpdateSegmentApplication(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE UPDATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	app := extractSegmentApplication(d)

	app.ID = d.Id()

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(updApp, app.SiteID)
	if err != nil {
		return err
	}

	setSegmentApplicationFields(d, updApp, client.TenantBaseDomain)

	return resourceReadSegmentApplication(d, m)
}

func resourceDeleteSegmentApplication(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE DELETE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	err := client.Applications.DeleteApplication(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")

	return resourceReadSegmentApplication(d, m)
}

func extractSegmentApplication(d *schema.ResourceData) *dto.Application {
	segmentSettings := extractSegmentSettings(d)
	return &dto.Application{
		Name:                 d.Get("name").(string),
		Icon:                 d.Get("icon").(string),
		SiteID:               d.Get("site_id").(string),
		Type:                 "segment",
		Visible:              d.Get("visible").(bool),
		NotificationsEnabled: d.Get("notification_enabled").(bool),
		Subdomain:            d.Get("subdomain").(string),
		ExternalAddress:      d.Get("external_address").(string),
		SegmentSettings:      segmentSettings,
	}
}

func setSegmentApplicationFields(d *schema.ResourceData, application *dto.Application, tenantBaseDomain string) {
	d.Set("name", application.Name)
	d.Set("icon", application.Icon)
	d.Set("type", application.Type)
	d.Set("visible", application.Visible)
	d.Set("notification_enabled", application.NotificationsEnabled)
	d.Set("internal_address", application.InternalAddress)
	d.Set("external_address", application.ExternalAddress)
	if application.SegmentSettings != nil {
		d.Set("segment_settings", flattenSegmentSettings(application.SegmentSettings))
	}
}

func flattenSegmentSettings(settings *dto.SegmentSettings) []interface{} {
	if settings == nil {
		return []interface{}{}
	}

	k := map[string]interface{}{
		"original_ip": settings.OriginalIP,
	}

	return []interface{}{k}
}

func extractSegmentSettings(d *schema.ResourceData) *dto.SegmentSettings {
	var segmentSettings *dto.SegmentSettings

	if v, ok := d.GetOk("segment_settings"); ok {
		for _, element := range v.([]interface{}) {
			elem := element.(map[string]interface{})

			var originalIp string

			originalIp = elem["original_ip"].(string)

			segmentSettings = &dto.SegmentSettings{
				OriginalIP: originalIp,
			}
		}
	}

	return segmentSettings
}
