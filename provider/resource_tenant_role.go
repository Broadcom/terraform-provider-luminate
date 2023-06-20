package provider

import (
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func LuminateTenantRole() *schema.Resource {
	tenantRoleSchema := LuminateAssignRoleBaseSchema()
	return &schema.Resource{
		Schema: tenantRoleSchema,
		Create: resourceCreateTenantRole,
		Read:   resourceReadTenantRole,
		Delete: resourceDeleteRoleBindings,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateTenantRole(d *schema.ResourceData, m interface{}) error {
	log.Println("[Info] Creating Tenant Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client ")
	}

	baseFields, err := extractBaseFields(d)
	if err != nil {
		return errors.Wrap(err, "extract tenant role base fields error")
	}

	identityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(baseFields.EntityIDPID)
	if err != nil {
		return errors.Wrapf(err, "Failed to lookup identity provider type for identity provider id %s", baseFields.EntityIDPID)
	}

	entity := dto.DirectoryEntity{
		IdentifierInProvider: baseFields.EntityIDInIDP,
		IdentityProviderId:   baseFields.EntityIDPID,
		EntityType:           baseFields.EntityType,
		IdentityProviderType: dto.ConvertIdentityProviderTypeToString(identityProviderType),
		DisplayName:          "displayName",
	}

	var tenantRole = dto.CreateRoleDTO{
		Role:     baseFields.RoleType,
		Entities: []dto.DirectoryEntity{entity},
	}

	roleBindings, err := client.RoleBindingsAPI.CreateTenantRoleBindings(&tenantRole)
	if err != nil {
		return err
	}

	d.SetId(roleBindings[0].ID)
	return nil
}

func resourceReadTenantRole(d *schema.ResourceData, m interface{}) error {
	log.Println("[Info] Reading Tenant Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	roleID := d.Id()
	roleType := d.Get("role_type").(string)
	if !utils.ValidateTenantRole(roleType) {
		return errors.New("invalid role type")
	}
	entityID := d.Get("entity_id").(string)
	role, err := client.RoleBindingsAPI.ReadRoleBindings(roleID, roleType, entityID, "", "")
	if err != nil {
		return errors.Wrap(err, "read tenant role failure")
	}

	d.SetId(role.ID)
	return nil
}
