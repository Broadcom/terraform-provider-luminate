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

func LuminateTcpAccessPolicy() *schema.Resource {
	tcpSchema := LuminateAccessPolicyBaseSchema()

	tcpSchema["allow_temporary_token"] = &schema.Schema{
		Type:         schema.TypeBool,
		Description:  "Indication whether authentication using a temporary token is allowed.",
		Optional:     true,
		Default:      true,
		ValidateFunc: utils.ValidateBool,
	}

	tcpSchema["allow_public_key"] = &schema.Schema{
		Type:         schema.TypeBool,
		Description:  "Indication whether authentication using long term secret is allowed.",
		Optional:     true,
		Default:      false,
		ValidateFunc: utils.ValidateBool,
	}

	return &schema.Resource{
		Schema:        tcpSchema,
		CreateContext: resourceCreateTcpAccessPolicy,
		ReadContext:   resourceReadTcpAccessPolicy,
		UpdateContext: resourceUpdateTcpAccessPolicy,
		DeleteContext: resourceDeleteAccessPolicy,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateTcpAccessPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	accessPolicy := extractTcpAccessPolicy(d)
	for i, _ := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			return diag.FromErr(errors.Wrapf(err, "Failed to lookup identity provider type for identity provider id %s", accessPolicy.DirectoryEntities[i].IdentityProviderId))
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)
	}

	createdAccessPolicy, err := client.AccessPolicies.CreateAccessPolicy(accessPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	setTcpAccessPolicyFields(d, createdAccessPolicy)

	resourceReadTcpAccessPolicy(ctx, d, m)
	return diagnostics
}

func resourceReadTcpAccessPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	setTcpAccessPolicyFields(d, accessPolicy)

	return diagnostics
}

func resourceUpdateTcpAccessPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	accessPolicy := extractTcpAccessPolicy(d)
	for i, _ := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			return diag.FromErr(errors.Wrapf(err, "Failed to lookup identity provider type for identity provider id %s", accessPolicy.DirectoryEntities[i].IdentityProviderId))
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)
	}
	accessPolicy.Id = d.Id()

	accessPolicy, err := client.AccessPolicies.UpdateAccessPolicy(accessPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	setTcpAccessPolicyFields(d, accessPolicy)

	resourceReadTcpAccessPolicy(ctx, d, m)
	return diagnostics
}

func setTcpAccessPolicyFields(d *schema.ResourceData, accessPolicy *dto.AccessPolicy) {
	setAccessPolicyBaseFields(d, accessPolicy)

	d.Set("allow_temporary_token", accessPolicy.TcpSettings.AcceptTemporaryToken)
	d.Set("allow_public_key", accessPolicy.TcpSettings.AcceptCertificate)
}

func extractTcpAccessPolicy(d *schema.ResourceData) *dto.AccessPolicy {
	accessPolicy := extractAccessPolicyBaseFields(d)

	accessPolicy.TargetProtocol = "TCP"
	accessPolicy.TcpSettings = &dto.PolicyTcpSettings{
		AcceptTemporaryToken: d.Get("allow_temporary_token").(bool),
		AcceptCertificate:    d.Get("allow_public_key").(bool),
	}

	return accessPolicy
}
