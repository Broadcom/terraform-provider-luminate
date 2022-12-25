package provider

import (
	"errors"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

func LuminateSshGwApplication() *schema.Resource {
	sshGwSchema := CommonApplicationSchema()

	sshGwSchema["integration_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "integration id used to setup the ssh gw application",
	}

	sshGwSchema["tags"] = &schema.Schema{
		Type:        schema.TypeMap,
		Required:    true,
		Description: "a map of tags used to determine which machines is included as part of this ssh gw",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	sshGwSchema["vpc"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		Required:    true,
		Description: "A list of vpc definitions used to determine the target group to include as part of the ssh gw application",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"vpc_id": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "the vpc id of the vpc containing target machines",
					ValidateFunc: utils.ValidateString,
				},
				"region": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "the region containing the target machines",
					ValidateFunc: utils.ValidateString,
				},
				"cidr_block": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "the cidr block of the machines to include",
					ValidateFunc: utils.ValidateString,
				},
			},
		},
	}

	sshGwSchema["segment_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		Schema: sshGwSchema,
		Create: resourceCreateSshGwApplication,
		Read:   resourceReadSshGwApplication,
		Update: resourceUpdateSshGwApplication,
		Delete: resourceDeleteSshGwApplication,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateSshGwApplication(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] LUMINATE CREATE SSH-GW APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	newApp := extractSshGwApplicationFields(d)

	app, err := client.Applications.CreateApplication(newApp)
	if err != nil {
		return err
	}

	d.SetId(app.ID)

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return err
	}

	setSshGwApplicationFields(d, app)

	return resourceReadSshGwApplication(d, m)
}

func resourceReadSshGwApplication(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] LUMINATE READ SSH-GW APP")

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

	d.SetId(app.ID)
	app.SiteID = d.Get("site_id").(string)
	setSshGwApplicationFields(d, app)

	return nil
}

func resourceUpdateSshGwApplication(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE UPDATE SSH-GW APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	app := extractSshGwApplicationFields(d)

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(app, app.SiteID)
	if err != nil {
		return err
	}

	setSshGwApplicationFields(d, updApp)

	return resourceReadSshGwApplication(d, m)
}

func resourceDeleteSshGwApplication(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE DELETE SSH-GW APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	err := client.Applications.DeleteApplication(d.Id())
	if err != nil {
		return err
	}

	return resourceReadSshGwApplication(d, m)
}

func setSshGwApplicationFields(d *schema.ResourceData, application *dto.Application) {
	d.Set("name", application.Name)
	d.Set("icon", application.Icon)
	d.Set("type", application.Type)
	d.Set("visible", application.Visible)
	d.Set("notification_enabled", application.NotificationsEnabled)
	d.Set("internal_address", application.InternalAddress)
	d.Set("external_address", application.ExternalAddress)
	d.Set("subdomain", application.Subdomain)
	d.Set("luminate_address", application.LuminateAddress)
	d.Set("segment_id", application.CloudIntegrationData.SegmentId)
	d.Set("tags", application.CloudIntegrationData.Tags)
	d.Set("vpc", flattenVpcs(application.CloudIntegrationData.Vpcs))
}

func flattenVpcs(vpcs []dto.Vpc) []interface{} {

	var flattenedVpcs []interface{}

	for _, c := range vpcs {
		flattenedVpc := map[string]interface{}{
			"IntegrationId": c.IntegrationId,
			"Region":        c.Region,
			"Vpc":           c.Vpc,
			"CidrBlock":     c.CidrBlock,
		}

		flattenedVpcs = append(flattenedVpcs, flattenedVpc)
	}
	return flattenedVpcs
}

func extractSshGwApplicationFields(d *schema.ResourceData) *dto.Application {
	tags := extractTagsField(d)
	integrationId := d.Get("integration_id").(string)
	vpcs := extractSshGwVpc(d, integrationId)

	var segmentId string
	if segmentIdInterface, ok := d.GetOkExists("segment_id"); ok {
		segmentId = segmentIdInterface.(string)
	}

	cloudIntegrationData := &dto.CloudIntegrationData{
		Tags:      tags,
		Vpcs:      vpcs,
		SegmentId: segmentId,
	}

	return &dto.Application{
		ID:                   d.Id(),
		Name:                 d.Get("name").(string),
		Icon:                 d.Get("icon").(string),
		SiteID:               d.Get("site_id").(string),
		Type:                 "sshgw",
		Visible:              d.Get("visible").(bool),
		NotificationsEnabled: d.Get("notification_enabled").(bool),
		ExternalAddress:      d.Get("external_address").(string),
		Subdomain:            d.Get("subdomain").(string),
		CloudIntegrationData: cloudIntegrationData,
	}
}

func extractTagsField(d *schema.ResourceData) map[string]string {
	tagsConfig := d.Get("tags").(map[string]interface{})

	tags := map[string]string{}

	for key, value := range tagsConfig {
		tags[key] = value.(string)
	}

	return tags
}

func extractSshGwVpc(d *schema.ResourceData, integrationId string) []dto.Vpc {
	var vpcs []dto.Vpc

	vpcsConfigs, ok := d.Get("vpc").([]interface{})
	if !ok {
		return vpcs
	}

	for _, v := range vpcsConfigs {
		vt := v.(map[string]interface{})

		vpc := dto.Vpc{
			IntegrationId: integrationId,
			Region:        vt["region"].(string),
			CidrBlock:     vt["cidr_block"].(string),
			Vpc:           vt["vpc_id"].(string),
		}

		vpcs = append(vpcs, vpc)
	}

	return vpcs
}
