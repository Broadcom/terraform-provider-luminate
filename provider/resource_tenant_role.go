package provider

import (
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	entity, err := getEntityByRoleBindings(client, baseFields)
	if err != nil {
		return err
	}
	roleType := d.Get("role_type").(string)
	if !utils.ValidateTenantRole(roleType) {
		return errors.New("invalid role type")
	}
	var tenantRole = dto.CreateRoleDTO{
		Role:     roleType,
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
