// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

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

func LuminateTenantRole() *schema.Resource {
	tenantRoleSchema := LuminateAssignRoleBaseSchema()
	return &schema.Resource{
		Schema:        tenantRoleSchema,
		CreateContext: resourceCreateTenantRole,
		ReadContext:   resourceReadTenantRole,
		DeleteContext: resourceDeleteRoleBindings,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateTenantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Creating Tenant Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client "))
	}

	baseFields, err := extractBaseFields(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "extract tenant role base fields error"))
	}

	entity, err := getEntityByRoleBindings(client, baseFields)
	if err != nil {
		return diag.FromErr(err)
	}
	roleType := d.Get("role_type").(string)
	if !utils.ValidateTenantRole(roleType) {
		return diag.FromErr(errors.New("invalid role type"))
	}
	var tenantRole = dto.CreateRoleDTO{
		Role:     roleType,
		Entities: []dto.DirectoryEntity{entity},
	}

	roleBindings, err := client.RoleBindingsAPI.CreateTenantRoleBindings(&tenantRole)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(roleBindings[0].ID)
	return resourceReadTenantRole(ctx, d, m)
}

func resourceReadTenantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Reading Tenant Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	roleID := d.Id()
	roleType := d.Get("role_type").(string)
	if !utils.ValidateTenantRole(roleType) {
		return diag.FromErr(errors.New("invalid role type"))
	}
	entityID := d.Get("entity_id").(string)
	role, err := client.RoleBindingsAPI.ReadRoleBindings(roleID, roleType, entityID, "", "")
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "read tenant role failure"))
	}

	d.SetId(role.ID)
	return nil
}
