package provider

import (
	"bitbucket.org/accezz-io/terraform-provider-symcsc/service"
	"bitbucket.org/accezz-io/terraform-provider-symcsc/service/dto"
	"bitbucket.org/accezz-io/terraform-provider-symcsc/utils"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func LuminateWebApplication() *schema.Resource {
	webAppSchema := CommonApplicationSchema()

	webAppSchema["internal_address"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "Internal address of the application, accessable by connector",
	}

	webAppSchema["custom_root_path"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "Requests coming into the external address root path '/', will be redirected to this custom path instead.",
	}
	webAppSchema["health_url"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "Health check path. The URI is relative to the external address.",
	}
	webAppSchema["health_method"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validateHealthMethod,
		Description:  "HTTP method to validate application health.",
	}
	webAppSchema["default_content_rewrite_rules_enabled"] = &schema.Schema{
		Type:         schema.TypeBool,
		Optional:     true,
		Default:      true,
		ValidateFunc: utils.ValidateBool,
		Description:  "Indicates whether to enable automatic translation of all occurrences of the application internal address to its external address on most prominent content types and relevant headers.",
	}
	webAppSchema["default_header_rewrite_rules_enabled"] = &schema.Schema{
		Type:         schema.TypeBool,
		Optional:     true,
		Default:      true,
		ValidateFunc: utils.ValidateBool,
		Description:  "Indicates whether to enable automatic translation of all occurrences of the application internal address to its external address on relevant headers.",
	}
	webAppSchema["use_external_address_for_host_and_sni"] = &schema.Schema{
		Type:         schema.TypeBool,
		Optional:     true,
		Default:      false,
		ValidateFunc: utils.ValidateBool,
		Description:  "Indicates whether to use external address for host header and SNI.",
	}
	webAppSchema["linked_applications"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "This property should be set in a scenario where the defined application contains resources that reference additional web applications by their internal domain name.",
	}

	webAppSchema["header_customization"] = &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "Custom headers key:value pairs to be added to all requests.",
	}

	return &schema.Resource{
		Schema: webAppSchema,
		Create: resourceCreateWebApplication,
		Read:   resourceReadWebApplication,
		Update: resourceUpdateWebApplication,
		Delete: resourceDeleteWebApplication,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateWebApplication(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] LUMINATE CREATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	newApp := extractWebApplication(d)

	app, err := client.Applications.CreateApplication(newApp)

	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return err
	}

	d.SetId(app.ID)
	setWebApplicationFields(d, app)

	return resourceReadWebApplication(d, m)
}

func resourceReadWebApplication(d *schema.ResourceData, m interface{}) error {

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
	setWebApplicationFields(d, app)

	return nil
}

func resourceUpdateWebApplication(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE UPDATE APP")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	app := extractWebApplication(d)

	app.ID = d.Id()

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return err
	}

	err = client.Applications.BindApplicationToSite(updApp, app.SiteID)
	if err != nil {
		return err
	}

	setWebApplicationFields(d, updApp)

	return resourceReadWebApplication(d, m)
}

func resourceDeleteWebApplication(d *schema.ResourceData, m interface{}) error {
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

	return resourceReadWebApplication(d, m)
}

func validateHealthMethod(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	cType, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type to be string"))
		return warns, errs
	}

	validTypes := []string{
		"head",
		"get",
	}

	if !utils.StringInSlice(validTypes, cType) {
		errs = append(errs, fmt.Errorf("health_method type must be one of %v", validTypes))
	}
	return warns, errs
}

func extractWebApplication(d *schema.ResourceData) *dto.Application {
	return &dto.Application{
		Name:                              d.Get("name").(string),
		Icon:                              d.Get("icon").(string),
		SiteID:                            d.Get("site_id").(string),
		Type:                              "web",
		Visible:                           d.Get("visible").(bool),
		NotificationsEnabled:              d.Get("notification_enabled").(bool),
		InternalAddress:                   d.Get("internal_address").(string),
		ExternalAddress:                   d.Get("external_address").(string),
		Subdomain:                         d.Get("subdomain").(string),
		CustomExternalAddress:             d.Get("custom_external_address").(string),
		CustomRootPath:                    d.Get("custom_root_path").(string),
		HealthURL:                         d.Get("health_url").(string),
		HealthMethod:                      d.Get("health_method").(string),
		DefaultContentRewriteRulesEnabled: d.Get("default_content_rewrite_rules_enabled").(bool),
		DefaultHeaderRewriteRulesEnabled:  d.Get("default_header_rewrite_rules_enabled").(bool),
		UseExternalAddressForHostAndSni:   d.Get("use_external_address_for_host_and_sni").(bool),
		LinkedApplications:                expandStringList(d.Get("linked_applications").([]interface{})),
		HeaderCustomization:               d.Get("header_customization").(map[string]interface{}),
	}
}

func setWebApplicationFields(d *schema.ResourceData, application *dto.Application) {
	d.Set("name", application.Name)
	d.Set("icon", application.Icon)
	d.Set("type", application.Type)
	d.Set("visible", application.Visible)
	d.Set("notification_enabled", application.NotificationsEnabled)
	d.Set("internal_address", application.InternalAddress)
	d.Set("external_address", application.ExternalAddress)
	d.Set("subdomain", application.Subdomain)
	d.Set("custom_external_address", application.CustomExternalAddress)
	d.Set("custom_root_path", application.CustomRootPath)
	d.Set("luminate_address", application.LuminateAddress)

	d.Set("health_url", application.HealthURL)
	d.Set("health_method", application.HealthMethod)
	d.Set("default_content_rewrite_rules_enabled", application.DefaultContentRewriteRulesEnabled)
	d.Set("default_header_rewrite_rules_enabled", application.DefaultHeaderRewriteRulesEnabled)
	d.Set("use_external_address_for_host_and_sni", application.UseExternalAddressForHostAndSni)
	d.Set("linked_applications", application.LinkedApplications)
	d.Set("header_customization", application.HeaderCustomization)
}
