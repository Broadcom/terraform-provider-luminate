// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/pkg/errors"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
)

func LuminateRdpAccessPolicy() *schema.Resource {
	rdpSchema := LuminateAccessPolicyBaseSchema()

	rdpSchema["allow_long_term_password"] = &schema.Schema{
		Type:         schema.TypeBool,
		Description:  "Indicates whether to allow long term password.",
		Optional:     true,
		Default:      false,
		ValidateFunc: utils.ValidateBool,
	}

	rdpSchema["target_protocol_subtype"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      string(sdk.NATIVE_PolicyTargetProtocolSubType),
		ValidateFunc: validateRdpTargetProtocolSubType,
		Description:  "rdp policy target protocol sub type",
	}

	rdpSchema["web_rdp_settings"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "Web RDP settings.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"disable_copy": {
					Type:        schema.TypeBool,
					Description: "Indicates whether to disable copy.",
					Default:     false,
					Optional:    true,
				},
				"disable_paste": {
					Type:        schema.TypeBool,
					Description: "Indicates whether to disable paste.",
					Default:     true,
					Optional:    true,
				},
			},
		},
	}

	return &schema.Resource{
		Schema:        rdpSchema,
		CreateContext: resourceCreateRdpAccessPolicy,
		ReadContext:   resourceReadRdpAccessPolicy,
		UpdateContext: resourceUpdateRdpAccessPolicy,
		DeleteContext: resourceDeleteAccessPolicy,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func validateRdpTargetProtocolSubType(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	cType, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type to be string"))
		return warns, errs
	}

	validTypes := []string{
		string(sdk.NATIVE_PolicyTargetProtocolSubType),
		string(sdk.BROWSER_PolicyTargetProtocolSubType),
	}

	if !utils.StringInSlice(validTypes, cType) {
		errs = append(errs, fmt.Errorf("target_protocol_subtype must be one of %v", validTypes))
	}
	return warns, errs
}

func resourceCreateRdpAccessPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	accessPolicy := extractRdpAccessPolicy(d)
	for i := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			err = errors.Wrapf(err, "Failed to lookup identity provider type for identity provider id %s", accessPolicy.DirectoryEntities[i].IdentityProviderId)
			return diag.FromErr(err)
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)
	}

	createdAccessPolicy, err := client.AccessPolicies.CreateAccessPolicy(accessPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	err = setRdpAccessPolicyFields(d, createdAccessPolicy)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "Failed to set access policy field"))
	}

	resourceReadRdpAccessPolicy(ctx, d, m)
	return diagnostics
}

func resourceReadRdpAccessPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	accessPolicy, err := client.AccessPolicies.GetAccessPolicy(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if accessPolicy == nil {
		d.SetId("")
		return nil
	}

	err = setRdpAccessPolicyFields(d, accessPolicy)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "Failed to set access policy field"))
	}

	return diagnostics
}

func resourceUpdateRdpAccessPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	accessPolicy := extractRdpAccessPolicy(d)
	for i := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			return diag.FromErr(errors.Wrapf(err, "Failed to lookup identity provider type for identity provider id %s", accessPolicy.DirectoryEntities[i].IdentityProviderId))
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)
	}
	accessPolicy.Id = d.Id()

	updatedAccessPolicy, err := client.AccessPolicies.UpdateAccessPolicy(accessPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	err = setRdpAccessPolicyFields(d, updatedAccessPolicy)
	if err != nil {
		return diag.FromErr(errors.Wrapf(err, "Failed to set access policy field"))
	}

	resourceReadRdpAccessPolicy(ctx, d, m)
	return diagnostics
}

func setRdpAccessPolicyFields(d *schema.ResourceData, accessPolicy *dto.AccessPolicy) error {
	setAccessPolicyBaseFields(d, accessPolicy)
	if err := d.Set("allow_long_term_password", accessPolicy.RdpSettings.LongTermPassword); err != nil {
		return err
	}
	if accessPolicy.RdpSettings.WebRdpSettings != nil {
		webSettings := map[string]interface{}{
			"disable_copy":  accessPolicy.RdpSettings.WebRdpSettings.DisableCopy,
			"disable_paste": accessPolicy.RdpSettings.WebRdpSettings.DisablePaste,
		}
		if err := d.Set("web_rdp_settings", []interface{}{webSettings}); err != nil {
			return err
		}
	}

	return nil
}

func extractRdpAccessPolicy(d *schema.ResourceData) *dto.AccessPolicy {
	accessPolicy := extractAccessPolicyBaseFields(d)

	longTermPassword := d.Get("allow_long_term_password").(bool)
	targetProtocolSubtype := d.Get("target_protocol_subtype").(string)

	accessPolicy.TargetProtocol = "RDP"
	accessPolicy.TargetProtocolSubtype = targetProtocolSubtype
	accessPolicy.RdpSettings = &dto.PolicyRdpSettings{
		LongTermPassword: longTermPassword,
	}

	if targetProtocolSubtype == string(sdk.BROWSER_PolicyTargetProtocolSubType) {
		// default WebRDP settings
		webRdpSettings := &dto.PolicyWebRdpSettings{
			DisableCopy:  false,
			DisablePaste: true,
		}
		if v, ok := d.GetOk("web_rdp_settings"); ok {
			settingsList := v.([]interface{})
			if len(settingsList) > 0 && settingsList[0] != nil {
				settingsMap := settingsList[0].(map[string]interface{})
				webRdpSettings.DisableCopy = settingsMap["disable_copy"].(bool)
				webRdpSettings.DisablePaste = settingsMap["disable_paste"].(bool)
			}
		}
		accessPolicy.RdpSettings.WebRdpSettings = webRdpSettings
	}

	return accessPolicy
}
