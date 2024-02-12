package provider

import (
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"log"
	"strings"
)

func LuminateWebAccessPolicy() *schema.Resource {
	webSchema := LuminateAccessPolicyBaseSchema()
	// web application is the only one that have mfa verification on access policy
	webSchema["validators"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"mfa": {
					Type:         schema.TypeBool,
					Optional:     true,
					Default:      false,
					Description:  "Indicate whatever to mfa verification validation.",
					ValidateFunc: utils.ValidateBool,
				},
			},
		},
	}

	return &schema.Resource{
		Schema:        webSchema,
		CreateContext: resourceCreateWebAccessPolicy,
		ReadContext:   resourceReadWebAccessPolicy,
		UpdateContext: resourceUpdateWebAccessPolicy,
		DeleteContext: resourceDeleteAccessPolicy,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateWebAccessPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Creating web access policy")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	accessPolicy := extractWebAccessPolicy(d)

	for i := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			return diag.FromErr(errors.Wrapf(err, "Failed to lookup identity provider type for identity provider id %s", accessPolicy.DirectoryEntities[i].IdentityProviderId))
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)

		// Get Display Name for User/Group by ID
		var resolvedDisplayName string
		switch strings.ToLower(accessPolicy.DirectoryEntities[i].EntityType) {
		case "user":
			resolvedDisplayName, err = client.IdentityProviders.GetUserDisplayNameTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId, accessPolicy.DirectoryEntities[i].IdentifierInProvider)
		case "group":
			resolvedDisplayName, err = client.IdentityProviders.GetGroupDisplayNameTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId, accessPolicy.DirectoryEntities[i].IdentifierInProvider)
		default:
			return diag.FromErr(errors.Wrapf(err, "Failed to lookup displayName - unknown entity type \"%s\"", accessPolicy.DirectoryEntities[i].EntityType))
		}

		if err != nil {
			return diag.FromErr(errors.Wrapf(err, "Failed to lookup displayName for entity type %s with identifier id %s on Identity Provider ID %s", accessPolicy.DirectoryEntities[i].EntityType, accessPolicy.DirectoryEntities[i].IdentifierInProvider, accessPolicy.DirectoryEntities[i].IdentityProviderId))
		}
		accessPolicy.DirectoryEntities[i].DisplayName = resolvedDisplayName
	}

	createdAccessPolicy, err := client.AccessPolicies.CreateAccessPolicy(accessPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	setAccessPolicyBaseFields(d, createdAccessPolicy)
	resourceReadWebAccessPolicy(ctx, d, m)
	return diagnostics
}

func resourceReadWebAccessPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading web access policy")
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

	setAccessPolicyBaseFields(d, accessPolicy)

	return diagnostics
}

func resourceUpdateWebAccessPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating web access policy")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	accessPolicy := extractWebAccessPolicy(d)
	for i := range accessPolicy.DirectoryEntities {
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

	setAccessPolicyBaseFields(d, accessPolicy)

	resourceReadWebAccessPolicy(ctx, d, m)
	return diagnostics
}

func extractWebAccessPolicy(d *schema.ResourceData) *dto.AccessPolicy {
	accessPolicy := extractAccessPolicyBaseFields(d)
	accessPolicy.TargetProtocol = "HTTP"

	return accessPolicy
}
