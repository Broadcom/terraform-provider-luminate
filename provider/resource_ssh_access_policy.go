package provider

import (
	"errors"
	"fmt"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func LuminateSshAccessPolicy() *schema.Resource {
	sshSchema := LuminateAccessPolicyBaseSchema()

	sshSchema["accounts"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "SSH/Unix accounts with which IDP entities and/or Luminate local users can access the SSH Server",
		Required:    true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.NoZeroValues,
		},
	}

	sshSchema["use_auto_mapping"] = &schema.Schema{
		Type: schema.TypeBool,
		Description: "Determine the strategy for mapping IDP entities to SSH/Unix accounts, " +
			"and specifically indicate whether automatic mapping based on the logged-in IDP entity username is allowed." +
			" In case this propert is set to TRUE, " +
			"manually entered SSH accounts are ignored. This property is relevant for SSH applications only.",
		Optional:     true,
		Default:      false,
		ValidateFunc: utils.ValidateBool,
	}

	sshSchema["full_upn_auto_mapping"] = &schema.Schema{
		Type: schema.TypeBool,
		Description: "Determine the strategy for mapping IDP entities to SSH/Unix accounts. In case this " +
			"property is set to true, full UPN is used, otherwise username(the username is the " +
			"part before the @ of the user’s UPN)is used. This property applies only in case " +
			"autoMapping is set to true.",
		Optional:     true,
		Default:      false,
		ValidateFunc: utils.ValidateBool,
	}

	sshSchema["allow_agent_forwarding"] = &schema.Schema{
		Type: schema.TypeBool,
		Description: "Indicates whether SSH agent forwarding is allowed for a transparent secure access to all " +
			"corporate SSH Servers via this SSH application that acts a Bastion." +
			" This property is relevant for SSH applications only.",
		Optional:     true,
		Default:      false,
		ValidateFunc: utils.ValidateBool,
	}

	sshSchema["allow_temporary_token"] = &schema.Schema{
		Type:         schema.TypeBool,
		Description:  "Indication whether authentication using a temporary token is allowed.",
		Optional:     true,
		Default:      true,
		ValidateFunc: utils.ValidateBool,
	}

	sshSchema["allow_public_key"] = &schema.Schema{
		Type:         schema.TypeBool,
		Description:  "Indication whether authentication using long term secret is allowed.",
		Optional:     true,
		Default:      false,
		ValidateFunc: utils.ValidateBool,
	}

	return &schema.Resource{
		Schema: sshSchema,
		Create: resourceCreateSshAccessPolicy,
		Read:   resourceReadSshAccessPolicy,
		Update: resourceUpdateSshAccessPolicy,
		Delete: resourceDeleteAccessPolicy,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateSshAccessPolicy(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	accessPolicy := extractSshAccessPolicy(d)
	for i, _ := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			error := fmt.Sprintf("Failed to lookup identity provider type for identity provider id %s: %s", accessPolicy.DirectoryEntities[i].IdentityProviderId, err)
			return errors.New(error)
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)
	}

	createdAccessPolicy, err := client.AccessPolicies.CreateAccessPolicy(accessPolicy)
	if err != nil {
		return err
	}

	setSshAccessPolicyFields(d, createdAccessPolicy)

	return resourceReadSshAccessPolicy(d, m)
}

func resourceReadSshAccessPolicy(d *schema.ResourceData, m interface{}) error {
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

	setSshAccessPolicyFields(d, accessPolicy)

	return nil
}

func resourceUpdateSshAccessPolicy(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	accessPolicy := extractSshAccessPolicy(d)
	for i, _ := range accessPolicy.DirectoryEntities {
		resolvedIdentityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(accessPolicy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			error := fmt.Sprintf("Failed to lookup identity provider type for identity provider id %s: %s", accessPolicy.DirectoryEntities[i].IdentityProviderId, err)
			return errors.New(error)
		}
		accessPolicy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)
	}
	accessPolicy.Id = d.Id()

	accessPolicy, err := client.AccessPolicies.UpdateAccessPolicy(accessPolicy)
	if err != nil {
		return err
	}

	setSshAccessPolicyFields(d, accessPolicy)

	return resourceReadSshAccessPolicy(d, m)
}

func setSshAccessPolicyFields(d *schema.ResourceData, accessPolicy *dto.AccessPolicy) {
	setAccessPolicyBaseFields(d, accessPolicy)
	d.Set("accounts", accessPolicy.SshSettings.Accounts)
	d.Set("use_auto_mapping", accessPolicy.SshSettings.AutoMapping)
	d.Set("full_upn_auto_mapping", accessPolicy.SshSettings.FullUPNAutoMapping)
	d.Set("allow_agent_forwarding", accessPolicy.SshSettings.AgentForward)
	d.Set("allow_temporary_token", accessPolicy.SshSettings.AcceptTemporaryToken)
	d.Set("allow_public_key", accessPolicy.SshSettings.AcceptCertificate)
}

func extractSshAccessPolicy(d *schema.ResourceData) *dto.AccessPolicy {
	accessPolicy := extractAccessPolicyBaseFields(d)

	unixAccountsInterface := d.Get("accounts").([]interface{})
	unixAccounts := utils.ParseStringList(unixAccountsInterface)

	accessPolicy.TargetProtocol = "SSH"
	accessPolicy.SshSettings = &dto.PolicySshSettings{
		Accounts:             unixAccounts,
		AutoMapping:          d.Get("use_auto_mapping").(bool),
		FullUPNAutoMapping:   d.Get("full_upn_auto_mapping").(bool),
		AgentForward:         d.Get("allow_agent_forwarding").(bool),
		AcceptTemporaryToken: d.Get("allow_temporary_token").(bool),
		AcceptCertificate:    d.Get("allow_public_key").(bool),
	}

	return accessPolicy
}
