package provider

import (
	"errors"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

func LuminateWebAccessPolicy() *schema.Resource {
	webSchema := LuminateAccessPolicyBaseSchema()

	conditionsResource := webSchema["conditions"].Elem.(*schema.Resource)

	conditionsResource.Schema["managed_device"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Indicate whatever to restrict access to managed devices only",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"opswat": {
					Type:         schema.TypeBool,
					Optional:     true,
					Default:      false,
					Description:  "Indicate whatever to restrict access to Opswat MetaAccess",
					ValidateFunc: utils.ValidateBool,
				},
				"symantec_cloudsoc": {
					Type:         schema.TypeBool,
					Optional:     true,
					Default:      false,
					Description:  "Indicate whatever to restrict access to symantec cloudsoc",
					ValidateFunc: utils.ValidateBool,
				},
				"symantec_web_security_service": {
					Type:         schema.TypeBool,
					Optional:     true,
					Default:      false,
					Description:  "Indicate whatever to restrict access to symantec web security service",
					ValidateFunc: utils.ValidateBool,
				},
			},
		},
	}

	conditionsResource.Schema["unmanaged_device"] = &schema.Schema{
		Type:         schema.TypeBool,
		Optional:     true,
		Default:      false,
		Description:  "Indicate whatever to restrict access to unmanaged devices only",
		ValidateFunc: utils.ValidateBool,
	}

	return &schema.Resource{
		Schema: webSchema,
		Create: resourceCreateWebAccessPolicy,
		Read:   resourceReadWebAccessPolicy,
		Update: resourceUpdateWebAccessPolicy,
		Delete: resourceDeleteAccessPolicy,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateWebAccessPolicy(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	accessPolicy := extractWebAccessPolicy(d)

	for i, _ := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			error := fmt.Sprintf("Failed to lookup identity provider type for identity provider id %s: %s", accessPolicy.DirectoryEntities[i].IdentityProviderId, err)
			return errors.New(error)
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = resolvedIdentityProviderType

		// Get Display Name for User/Group by ID
		var resolvedDisplayName string
		switch accessPolicy.DirectoryEntities[i].EntityType {
		case "user", "User":
			resolvedDisplayName, err = client.IdentityProviders.GetUserDisplayNameTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId, accessPolicy.DirectoryEntities[i].IdentifierInProvider)
		case "group", "Group":
			resolvedDisplayName, err = client.IdentityProviders.GetGroupDisplayNameTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId, accessPolicy.DirectoryEntities[i].IdentifierInProvider)
		}

		if err != nil {
			error := fmt.Sprintf("Failed to lookup displayName for entity type %s with identifier id %s on Identity Provider ID %s: %v", accessPolicy.DirectoryEntities[i].EntityType, accessPolicy.DirectoryEntities[i].IdentifierInProvider, accessPolicy.DirectoryEntities[i].IdentityProviderId, err)
			return errors.New(error)
		}
		accessPolicy.DirectoryEntities[i].DisplayName = resolvedDisplayName
	}

	createdAccessPolicy, err := client.AccessPolicies.CreateAccessPolicy(accessPolicy)
	if err != nil {
		return err
	}

	setAccessPolicyBaseFields(d, createdAccessPolicy)

	return resourceReadWebAccessPolicy(d, m)
}

func resourceReadWebAccessPolicy(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	accessPolicy, err := client.AccessPolicies.GetAccessPolicy(d.Id())
	if err != nil {
		return err
	}

	if accessPolicy == nil {
		d.SetId("")
		return nil
	}

	setAccessPolicyBaseFields(d, accessPolicy)

	return nil
}

func resourceUpdateWebAccessPolicy(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	accessPolicy := extractWebAccessPolicy(d)
	for i, _ := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			error := fmt.Sprintf("Failed to lookup identity provider type for identity provider id %s: %s", accessPolicy.DirectoryEntities[i].IdentityProviderId, err)
			return errors.New(error)
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = resolvedIdentityProviderType
	}
	accessPolicy.Id = d.Id()

	accessPolicy, err := client.AccessPolicies.UpdateAccessPolicy(accessPolicy)
	if err != nil {
		return err
	}

	setAccessPolicyBaseFields(d, accessPolicy)

	return resourceReadWebAccessPolicy(d, m)
}

func extractWebAccessPolicy(d *schema.ResourceData) *dto.AccessPolicy {
	accessPolicy := extractAccessPolicyBaseFields(d)
	accessPolicy.TargetProtocol = "HTTP"

	return accessPolicy
}
