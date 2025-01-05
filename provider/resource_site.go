// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"kerberos": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Active Directory domain name you want to SSO with.",
							ValidateFunc: utils.ValidateString,
						},
						"kdc_address": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The hostname of the primary domain controller/domain controller closest to the connector.",
							ValidateFunc: utils.ValidateString,
						},
						"keytab_pair": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The absolute path of the keytab file",
							ValidateFunc: utils.ValidateString,
						},
					},
				},
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
		Name:       d.Get("name").(string),
		Region:     d.Get("region").(string),
		MuteHealth: d.Get("mute_health_notification").(bool),
		K8SVolume:  d.Get("kubernetes_persistent_volume_name").(string),
	}

	k, ok := d.Get("kerberos").(*schema.Set)

	if ok && len(k.List()) > 0 {
		kerb := k.List()[0].(map[string]interface{})

		site.Kerberos = &dto.SiteKerberosConfig{
			Domain:     kerb["domain"].(string),
			KDCAddress: kerb["kdc_address"].(string),
			KeytabPair: kerb["keytab_pair"].(string),
		}
	}
	return &site
}

func setSiteFields(d *schema.ResourceData, site *dto.Site) {
	d.Set("name", site.Name)
	d.Set("region", site.Region)
	d.Set("mute_health_notification", site.MuteHealth)
	d.Set("kubernetes_persistent_volume_name", site.K8SVolume)

	if site.Kerberos != nil && site.Kerberos.Domain != "" {
		d.Set("kerberos", flattenKerberosConfig(site.Kerberos))
	}
}

func flattenKerberosConfig(config *dto.SiteKerberosConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}
	k := map[string]interface{}{
		"domain":      config.Domain,
		"kdc_address": config.KDCAddress,
		"keytab_pair": config.KeytabPair,
	}
	return []interface{}{k}
}
