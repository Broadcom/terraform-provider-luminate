package provider

import (
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
		Schema: rdpSchema,
		Create: resourceCreateRdpAccessPolicy,
		Read:   resourceReadRdpAccessPolicy,
		Update: resourceUpdateRdpAccessPolicy,
		Delete: resourceDeleteAccessPolicy,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateRdpAccessPolicy(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	accessPolicy := extractRdpAccessPolicy(d)
	for i := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			return errors.Wrapf(err, "Failed to lookup identity provider type for identity provider id %s", accessPolicy.DirectoryEntities[i].IdentityProviderId)
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)
	}

	createdAccessPolicy, err := client.AccessPolicies.CreateAccessPolicy(accessPolicy)
	if err != nil {
		return err
	}

	err = setRdpAccessPolicyFields(d, createdAccessPolicy)
	if err != nil {
		return errors.Wrap(err, "Failed to set access policy field")
	}

	return resourceReadRdpAccessPolicy(d, m)
}

func resourceReadRdpAccessPolicy(d *schema.ResourceData, m interface{}) error {
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

	err = setRdpAccessPolicyFields(d, accessPolicy)
	if err != nil {
		return errors.Wrap(err, "Failed to set access policy field")
	}

	return nil
}

func resourceUpdateRdpAccessPolicy(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	accessPolicy := extractRdpAccessPolicy(d)
	for i := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			return errors.Wrapf(err, "Failed to lookup identity provider type for identity provider id %s", accessPolicy.DirectoryEntities[i].IdentityProviderId)
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)
	}
	accessPolicy.Id = d.Id()

	updatedAccessPolicy, err := client.AccessPolicies.UpdateAccessPolicy(accessPolicy)
	if err != nil {
		return err
	}

	err = setRdpAccessPolicyFields(d, updatedAccessPolicy)
	if err != nil {
		return errors.Wrapf(err, "Failed to set access policy field")
	}

	return resourceReadRdpAccessPolicy(d, m)
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
