package provider

import (
	"errors"
	"fmt"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
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
		Schema: tcpSchema,
		Create: resourceCreateTcpAccessPolicy,
		Read:   resourceReadTcpAccessPolicy,
		Update: resourceUpdateTcpAccessPolicy,
		Delete: resourceDeleteAccessPolicy,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateTcpAccessPolicy(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	accessPolicy := extractTcpAccessPolicy(d)
	for i, _ := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			error := fmt.Sprintf("Failed to lookup identity provider type for identity provider id %s: %s", accessPolicy.DirectoryEntities[i].IdentityProviderId, err)
			return errors.New(error)
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = resolvedIdentityProviderType
	}

	createdAccessPolicy, err := client.AccessPolicies.CreateAccessPolicy(accessPolicy)
	if err != nil {
		return err
	}

	setTcpAccessPolicyFields(d, createdAccessPolicy)

	return resourceReadTcpAccessPolicy(d, m)
}

func resourceReadTcpAccessPolicy(d *schema.ResourceData, m interface{}) error {
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

	setTcpAccessPolicyFields(d, accessPolicy)

	return nil
}

func resourceUpdateTcpAccessPolicy(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	accessPolicy := extractTcpAccessPolicy(d)
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

	setTcpAccessPolicyFields(d, accessPolicy)

	return resourceReadTcpAccessPolicy(d, m)
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
