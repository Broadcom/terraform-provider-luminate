package provider

import (
	"errors"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

func LuminateRDPApplication() *schema.Resource {
	rdpSchema := CommonApplicationSchema()

	rdpSchema["internal_address"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "Internal address of the application, accessible by connector",
	}

	return &schema.Resource{
		Schema: rdpSchema,
		Create: resourceCreateRDPApplication,
		Read:   resourceReadRDPApplication,
		Update: resourceUpdateRDPApplication,
		Delete: resourceDeleteRDPApplication,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateRDPApplication(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] LUMINATE CREATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	newApp := extractRDPApplicationFields(d)

	app, err := client.Applications.CreateApplication(newApp)
	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return err
	}

	d.SetId(app.ID)
	setSHHApplicationFields(d, app)

	return resourceReadRDPApplication(d, m)
}

func resourceReadRDPApplication(d *schema.ResourceData, m interface{}) error {

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

	d.SetId(app.ID)

	app.SiteID = d.Get("site_id").(string)
	setSHHApplicationFields(d, app)

	return nil
}

func resourceUpdateRDPApplication(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE UPDATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	app := extractRDPApplicationFields(d)

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(app, app.SiteID)
	if err != nil {
		return err
	}

	updApp.SiteID = app.SiteID
	setRDPApplicationFields(d, updApp)

	return resourceReadRDPApplication(d, m)
}

func resourceDeleteRDPApplication(d *schema.ResourceData, m interface{}) error {
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

	return resourceReadRDPApplication(d, m)
}

func setRDPApplicationFields(d *schema.ResourceData, application *dto.Application) {
	d.Set("name", application.Name)
	d.Set("icon", application.Icon)
	d.Set("type", application.Type)
	d.Set("visible", application.Visible)
	d.Set("site_id", application.SiteID)
	d.Set("notification_enabled", application.NotificationsEnabled)
	d.Set("internal_address", application.InternalAddress)
	d.Set("external_address", application.ExternalAddress)
	d.Set("subdomain", application.Subdomain)
	d.Set("custom_external_address", application.CustomExternalAddress)
	d.Set("luminate_address", application.LuminateAddress)
}

func extractRDPApplicationFields(d *schema.ResourceData) *dto.Application {
	return &dto.Application{
		ID:                    d.Id(),
		Name:                  d.Get("name").(string),
		Icon:                  d.Get("icon").(string),
		SiteID:                d.Get("site_id").(string),
		Type:                  "rdp",
		Visible:               d.Get("visible").(bool),
		NotificationsEnabled:  d.Get("notification_enabled").(bool),
		InternalAddress:       d.Get("internal_address").(string),
		ExternalAddress:       d.Get("external_address").(string),
		Subdomain:             d.Get("subdomain").(string),
		CustomExternalAddress: d.Get("custom_external_address").(string),
	}
}
