// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/pkg/errors"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	return d.Set("allow_long_term_password", accessPolicy.RdpSettings.LongTermPassword)
}

func extractRdpAccessPolicy(d *schema.ResourceData) *dto.AccessPolicy {
	accessPolicy := extractAccessPolicyBaseFields(d)

	longTermPassword := d.Get("allow_long_term_password").(bool)

	accessPolicy.TargetProtocol = "RDP"
	accessPolicy.RdpSettings = &dto.PolicyRdpSettings{
		LongTermPassword: longTermPassword,
	}

	return accessPolicy
}
