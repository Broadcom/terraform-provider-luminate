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
	"log"
)

const (
	FieldAuthenticationMode = "authentication_mode"
)

func LuminateSite() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Site name",
				ValidateFunc: utils.ValidateSiteName,
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Site connectivity region",
				ValidateFunc: utils.ValidateString,
			},
			FieldAuthenticationMode: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      dto.SiteAuthenticationModeConnector,
				Description:  "Site authentication mode",
				ValidateFunc: validateAuthenticationMode,
			},
			"mute_health_notification": {
				Type:         schema.TypeBool,
				Optional:     true,
				Description:  "Mute notifications if site is down",
				Default:      false,
				ValidateFunc: utils.ValidateBool,
			},
			"kubernetes_persistent_volume_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Kubernetes persistent volume name",
				ValidateFunc: utils.ValidateString,
			},
		},
		CreateContext: resourceCreateSite,
		ReadContext:   resourceReadSite,
		UpdateContext: resourceUpdateSite,
		DeleteContext: resourceDeleteSite,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func validateAuthenticationMode(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	mode, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type to be string"))
		return warns, errs
	}

	if !utils.StringInSlice(dto.ValidAuthenticationModes, mode) {
		errs = append(errs, fmt.Errorf("authentication mode must be one of %v", dto.ValidAuthenticationModes))
		return warns, errs
	}

	return warns, errs
}

func resourceCreateSite(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE SITE CREATE")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}
	site := extractSiteFields(d)

	newSite, err := client.Sites.CreateSite(site)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(site.ID)
	setSiteFields(d, newSite)

	return resourceReadSite(ctx, d, m)
}

func resourceReadSite(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE SITE READ")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	if client == nil {
		return diag.FromErr(errors.New("unable to initialize client"))
	}

	site, err := client.Sites.GetSiteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if site == nil {
		d.SetId("")
		return nil
	}

	setSiteFields(d, site)

	return diagnostics
}

func resourceUpdateSite(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE SITE UPDATE")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	site := extractSiteFields(d)

	s, err := client.Sites.UpdateSite(site, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	setSiteFields(d, s)

	return resourceReadSite(ctx, d, m)
}

func resourceDeleteSite(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE SITE DELETE")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	err := client.Sites.DeleteSite(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return resourceReadSite(ctx, d, m)
}

func extractSiteFields(d *schema.ResourceData) *dto.Site {
	site := dto.Site{
		Name:               d.Get("name").(string),
		Region:             d.Get("region").(string),
		AuthenticationMode: dto.SiteAuthenticationMode(d.Get(FieldAuthenticationMode).(string)),
		MuteHealth:         d.Get("mute_health_notification").(bool),
		K8SVolume:          d.Get("kubernetes_persistent_volume_name").(string),
	}

	return &site
}

func setSiteFields(d *schema.ResourceData, site *dto.Site) {
	d.Set("name", site.Name)
	d.Set("region", site.Region)
	d.Set(FieldAuthenticationMode, site.AuthenticationMode)
	d.Set("mute_health_notification", site.MuteHealth)
	d.Set("kubernetes_persistent_volume_name", site.K8SVolume)
}
