package provider

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"errors"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
)

func LuminateSegmentApplication() *schema.Resource {
	segmentAppSchema := CommonApplicationSchema()
	segmentAppSchema["sub_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: utils.ValidateString,
		Default:      string(sdk.SEGMENT_RANGE_ApplicationSubType),
		Description:  "The segment application sub type",
	}

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

	segmentAppSchema["multiple_segment_settings"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"original_ip": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "List of target IPs",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: utils.ValidateString,
					},
				},
			},
		},
	}

	return &schema.Resource{
		Schema:        segmentAppSchema,
		CreateContext: resourceCreateSegmentApplication,
		ReadContext:   resourceReadSegmentApplication,
		UpdateContext: resourceUpdateSegmentApplication,
		DeleteContext: resourceDeleteSegmentApplication,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateSegmentApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] LUMINATE CREATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}
	newApp := extractSegmentApplication(d)

	app, err := client.Applications.CreateApplication(newApp)

	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.ID)
	setSegmentApplicationFields(d, app, client.TenantBaseDomain)

	return resourceReadSegmentApplication(ctx, d, m)
}

func resourceReadSegmentApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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

	app.SiteID = d.Get("site_id").(string)
	setSegmentApplicationFields(d, app, client.TenantBaseDomain)

	return nil
}

func resourceUpdateSegmentApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE UPDATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}
	app := extractSegmentApplication(d)

	app.ID = d.Id()

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Applications.BindApplicationToSite(updApp, app.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	setSegmentApplicationFields(d, updApp, client.TenantBaseDomain)

	return resourceReadSegmentApplication(ctx, d, m)
}

func resourceDeleteSegmentApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE DELETE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	err := client.Applications.DeleteApplication(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceReadSegmentApplication(ctx, d, m)
}

func extractSegmentApplication(d *schema.ResourceData) *dto.Application {
	segmentSettings := extractSegmentSettings(d)
	multipleSegmentSettings := extractMultipleSegmentSettings(d)
	return &dto.Application{
		Name:                    d.Get("name").(string),
		Icon:                    d.Get("icon").(string),
		SiteID:                  d.Get("site_id").(string),
		Type:                    "segment",
		SubType:                 d.Get("sub_type").(string),
		Visible:                 d.Get("visible").(bool),
		NotificationsEnabled:    d.Get("notification_enabled").(bool),
		Subdomain:               d.Get("subdomain").(string),
		ExternalAddress:         d.Get("external_address").(string),
		SegmentSettings:         segmentSettings,
		MultipleSegmentSettings: multipleSegmentSettings,
	}
}

func setSegmentApplicationFields(d *schema.ResourceData, application *dto.Application, tenantBaseDomain string) {
	d.Set("name", application.Name)
	d.Set("icon", application.Icon)
	d.Set("type", application.Type)
	d.Set("sub_type", application.SubType)
	d.Set("visible", application.Visible)
	d.Set("notification_enabled", application.NotificationsEnabled)
	d.Set("external_address", application.ExternalAddress)
	if application.SegmentSettings != nil {
		d.Set("segment_settings", flattenSegmentSettings(application.SegmentSettings))
	}
	if application.MultipleSegmentSettings != nil {
		d.Set("multiple_segment_settings", flattenMultipleSegmentSettings(application.MultipleSegmentSettings))
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

func flattenMultipleSegmentSettings(settings []*dto.SegmentSettings) []interface{} {
	var originalIPs []interface{}
	for _, setting := range settings {
		originalIPs = append(originalIPs, setting.OriginalIP)
	}
	k := map[string]interface{}{
		"original_ip": originalIPs,
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

func extractMultipleSegmentSettings(d *schema.ResourceData) []*dto.SegmentSettings {
	var multipleSegmentSettings []*dto.SegmentSettings

	if v, ok := d.GetOk("multiple_segment_settings"); ok {
		for _, element := range v.([]interface{}) {
			elem := element.(map[string]interface{})

			ipsList := elem["original_ip"].([]interface{})

			for _, ip := range ipsList {
				multipleSegmentSettings = append(multipleSegmentSettings, &dto.SegmentSettings{
					OriginalIP: ip.(string)})
			}
		}
	}
	return multipleSegmentSettings
}
