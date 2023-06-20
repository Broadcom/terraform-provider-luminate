package provider

import (
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func LuminateSiteRole() *schema.Resource {
	siteSchema := LuminateAssignRoleBaseSchema()
	siteSchema["site_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Description:  "the site id to which this role assigned to.",
		Required:     true,
		ValidateFunc: utils.ValidateUuid,
		ForceNew:     true,
	}
	return &schema.Resource{
		Schema: siteSchema,
		Create: resourceCreateSiteRole,
		Read:   resourceReadSiteRole,
		Delete: resourceDeleteRoleBindings,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateSiteRole(d *schema.ResourceData, m interface{}) error {
	log.Println("[Info] Creating Site Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client ")
	}

	baseFields, err := extractBaseFields(d)
	if err != nil {
		return errors.Wrap(err, "extract tenant role base fields error")
	}
	if !utils.ValidateSiteRole(baseFields.RoleType) {
		return errors.New("invalid role type")
	}

	siteID := d.Get("site_id").(string)

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

	var tenantRole = dto.CreateSiteRoleDTO{
		CreateRoleDTO: dto.CreateRoleDTO{
			Role:     baseFields.RoleType,
			Entities: []dto.DirectoryEntity{entity},
		},
		SiteID: siteID,
	}

	roleBindings, err := client.RoleBindingsAPI.CreateSiteRoleBindings(&tenantRole)
	if err != nil {
		return err
	}

	d.SetId(roleBindings[0].ID)
	return nil
}

func resourceReadSiteRole(d *schema.ResourceData, m interface{}) error {
	log.Println("[Info] Reading Site Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	roleBindingsID := d.Id()
	roleType := d.Get("role_type").(string)
	entityID := d.Get("entity_id").(string)
	siteID := d.Get("site_id").(string)
	role, err := client.RoleBindingsAPI.ReadRoleBindings(roleBindingsID, roleType, entityID, "", siteID)
	if err != nil {
		return errors.Wrap(err, "read tenant role failure")
	}

	d.SetId(role.ID)
	return nil
}
