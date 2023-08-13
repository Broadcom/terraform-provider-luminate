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
)

func LuminateSiteRole() *schema.Resource {
	siteSchema := LuminateAssignRoleBaseSchema()
	siteSchema["site_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Description:  "The site id to which this role is assigned to.",
		Required:     true,
		ValidateFunc: utils.ValidateUuid,
		ForceNew:     true,
	}
	return &schema.Resource{
		Schema:        siteSchema,
		CreateContext: resourceCreateSiteRole,
		ReadContext:   resourceReadSiteRole,
		DeleteContext: resourceDeleteRoleBindings,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateSiteRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Creating Site Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client "))
	}

	baseFields, err := extractBaseFields(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "extract site role base fields error"))
	}
	if !utils.ValidateSiteRole(baseFields.RoleType) {
		return diag.FromErr(errors.New("invalid role type"))
	}

	siteID := d.Get("site_id").(string)

	entity, err := getEntityByRoleBindings(client, baseFields)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	d.SetId(roleBindings[0].ID)
	return resourceReadSiteRole(ctx, d, m)
}

func resourceReadSiteRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Reading Site Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	roleBindingsID := d.Id()
	roleType := d.Get("role_type").(string)
	if !utils.ValidateSiteRole(roleType) {
		return diag.FromErr(errors.New("invalid role type"))
	}
	entityID := d.Get("entity_id").(string)
	siteID := d.Get("site_id").(string)
	role, err := client.RoleBindingsAPI.ReadRoleBindings(roleBindingsID, roleType, entityID, "", siteID)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "read site role failure"))
	}

	d.SetId(role.ID)
	return nil
}
