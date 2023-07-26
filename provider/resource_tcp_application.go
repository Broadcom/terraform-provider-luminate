package provider

import (
	"errors"
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
			},
		},
	}

	return &schema.Resource{
		Schema: tcpSchema,
		Create: resourceCreateTCPApplication,
		Read:   resourceReadTCPApplication,
		Update: resourceUpdateTCPApplication,
		Delete: resourceDeleteTCPApplication,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateTCPApplication(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] LUMINATE CREATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	newApp := extractTCPApplicationFields(d)

	app, err := client.Applications.CreateApplication(newApp)
	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return err
	}

	d.SetId(app.ID)
	setTCPApplicationFields(d, app)

	return resourceReadTCPApplication(d, m)
}

func resourceReadTCPApplication(d *schema.ResourceData, m interface{}) error {

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
	setTCPApplicationFields(d, app)
	return nil
}

func resourceUpdateTCPApplication(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE UPDATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	app := extractTCPApplicationFields(d)

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(app, app.SiteID)
	if err != nil {
		return err
	}

	setTCPApplicationFields(d, updApp)

	return resourceReadTCPApplication(d, m)
}

func resourceDeleteTCPApplication(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE DELETE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	err := client.Applications.DeleteApplication(d.Id())
	if err != nil {
		return err
	}

	return resourceReadTCPApplication(d, m)
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

		target := dto.TCPTarget{
			Address: vt["address"].(string),
			Ports:   ports,
		}
		targets = append(targets, target)
	}
	return targets
}

func setTCPApplicationFields(d *schema.ResourceData, application *dto.Application) {
	d.Set("name", application.Name)
	d.Set("collection_id", application.CollectionID)
	d.Set("icon", application.Icon)
	d.Set("type", application.Type)
	d.Set("visible", application.Visible)
	d.Set("notification_enabled", application.NotificationsEnabled)
	d.Set("internal_address", TCPApplicationInternalAddress)
	d.Set("external_address", application.ExternalAddress)
	d.Set("subdomain", application.Subdomain)
	d.Set("luminate_address", application.LuminateAddress)
	log.Printf("[DEBUG] Settings TCP Targets")

	d.Set("target", flattenTCPTargets(application.Targets))
}

func flattenTCPTargets(targets []dto.TCPTarget) []interface{} {

	var flat []interface{}

	for _, c := range targets {
		t := map[string]interface{}{
			"address": c.Address,
			"ports":   c.Ports,
		}
		flat = append(flat, t)
	}
	return flat
}
