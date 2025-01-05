// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strings"

	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
)

func LuminateWebApplication() *schema.Resource {
	webAppSchema := CommonApplicationSchema()

	webAppSchema["sub_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      string(sdk.HTTP_LUMINATE_DOMAIN_ApplicationSubType),
		ValidateFunc: validateSubType,
		Description:  "Web application sub type",
	}
	webAppSchema["internal_address"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "Internal address of the application, accessible by connector",
	}
	webAppSchema["custom_root_path"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "Requests coming into the external address root path '/', will be redirected to this custom path instead.",
	}
	webAppSchema["custom_external_address"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "The application custom DNS address that exposes the application.",
	}

	webAppSchema["custom_domain_record_name"] = &schema.Schema{
		Type:         schema.TypeString,
		Computed:     true,
		Optional:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "Value of custom domain record",
	}

	webAppSchema["custom_domain_record_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Computed:     true,
		Optional:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "Type of custom domain record",
	}

	webAppSchema["wildcard_certificate"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "The certificate for the custom domain of the wildcard application.",
	}
	webAppSchema["wildcard_private_key"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "The private key of the certificate for the custom domain of the wildcard application.",
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
		Type:     schema.TypeSet,
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
		Schema:        webAppSchema,
		CreateContext: resourceCreateWebApplication,
		ReadContext:   resourceReadWebApplication,
		UpdateContext: resourceUpdateWebApplication,
		DeleteContext: resourceDeleteWebApplication,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateWebApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log.Printf("[DEBUG] LUMINATE CREATE APP")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}
	newApp := extractWebApplication(d)

	app, err := client.Applications.CreateApplication(newApp)

	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Applications.BindApplicationToSite(app, newApp.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(app.ID)
	setWebApplicationFields(d, app, client.TenantBaseDomain)
	resourceReadWebApplication(ctx, d, m)
	return diagnostics
}

func resourceReadWebApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	app.SiteID = d.Get("site_id").(string)
	setWebApplicationFields(d, app, client.TenantBaseDomain)
	log.Printf("[DEBUG] LUMINATE READ APP DONE")
	return diagnostics
}

func resourceUpdateWebApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE UPDATE APP")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}
	app := extractWebApplication(d)

	app.ID = d.Id()

	updApp, err := client.Applications.UpdateApplication(app)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Applications.BindApplicationToSite(updApp, app.SiteID)
	if err != nil {
		return diag.FromErr(err)
	}

	setWebApplicationFields(d, updApp, client.TenantBaseDomain)

	resourceReadWebApplication(ctx, d, m)
	return diagnostics
}

func resourceDeleteWebApplication(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	resourceReadWebApplication(ctx, d, m)
	return diagnostics
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

func validateSubType(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	cType, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type to be string"))
		return warns, errs
	}

	validTypes := []string{
		string(sdk.HTTP_LUMINATE_DOMAIN_ApplicationSubType),
		string(sdk.HTTP_CUSTOM_DOMAIN_ApplicationSubType),
		string(sdk.HTTP_WILDCARD_DOMAIN_ApplicationSubType),
		string(sdk.SINGLE_MACHINE_ApplicationSubType),
		string(sdk.MULTIPLE_MACHINES_ApplicationSubType),
	}

	if !utils.StringInSlice(validTypes, cType) {
		errs = append(errs, fmt.Errorf("sub_type must be one of %v", validTypes))
	}
	return warns, errs
}

func extractWebApplication(d *schema.ResourceData) *dto.Application {
	return &dto.Application{
		Name:                              d.Get("name").(string),
		Icon:                              d.Get("icon").(string),
		SiteID:                            d.Get("site_id").(string),
		CollectionID:                      d.Get("collection_id").(string),
		Type:                              "web",
		SubType:                           d.Get("sub_type").(string),
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
		LinkedApplications:                expandStringList(d.Get("linked_applications").(*schema.Set)),
		HeaderCustomization:               d.Get("header_customization").(map[string]interface{}),
		WildcardCertificate:               d.Get("wildcard_certificate").(string),
		WildcardPrivateKey:                d.Get("wildcard_private_key").(string),
	}
}

func setWebApplicationFields(d *schema.ResourceData, application *dto.Application, tenantBaseDomain string) {
	SetBaseApplicationFields(d, application)
	d.Set("sub_type", application.SubType)
	d.Set("internal_address", application.InternalAddress)
	d.Set("custom_external_address", application.CustomExternalAddress)
	d.Set("collection_id", application.CollectionID)

	if application.CustomExternalAddress != "" {
		trimmedAddress := strings.ReplaceAll(application.CustomExternalAddress, "https://", "")
		recordName := fmt.Sprintf("%s.%s", utils.StringMD5(trimmedAddress), tenantBaseDomain)
		d.Set("custom_domain_record_name", recordName)
		d.Set("custom_domain_record_type", "CNAME")
	}

	d.Set("custom_root_path", application.CustomRootPath)
	d.Set("luminate_address", application.LuminateAddress)
	d.Set("wildcard_certificate", application.WildcardCertificate)
	d.Set("wildcard_private_key", application.WildcardPrivateKey)
	d.Set("health_url", application.HealthURL)
	d.Set("health_method", application.HealthMethod)
	d.Set("default_content_rewrite_rules_enabled", application.DefaultContentRewriteRulesEnabled)
	d.Set("default_header_rewrite_rules_enabled", application.DefaultHeaderRewriteRulesEnabled)
	d.Set("use_external_address_for_host_and_sni", application.UseExternalAddressForHostAndSni)
	d.Set("linked_applications", application.LinkedApplications)
	d.Set("header_customization", application.HeaderCustomization)

}
