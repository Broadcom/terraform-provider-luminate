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

func LuminateCollectionRole() *schema.Resource {
	collectionSchema := LuminateAssignRoleBaseSchema()
	collectionSchema["collection_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Description:  "The collection id to which this role will be assigned to.",
		Required:     true,
		ValidateFunc: utils.ValidateUuid,
		ForceNew:     true,
	}
	return &schema.Resource{
		Schema:        collectionSchema,
		CreateContext: resourceCreateCollectionRole,
		ReadContext:   resourceReadCollectionRole,
		DeleteContext: resourceDeleteRoleBindings,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateCollectionRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Creating Collection Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client "))
	}

	baseFields, err := extractBaseFields(d)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "extract collection role base fields error"))
	}
	if !utils.ValidateCollectionRole(baseFields.RoleType) {
		return diag.FromErr(errors.New("invalid role type"))
	}

	collectionID := d.Get("collection_id").(string)

	entity, err := getEntityByRoleBindings(client, baseFields)
	if err != nil {
		return diag.FromErr(err)
	}

	var collectionRole = dto.CreateCollectionRoleDTO{
		CreateRoleDTO: dto.CreateRoleDTO{
			Role:     baseFields.RoleType,
			Entities: []dto.DirectoryEntity{entity},
		},
		CollectionID: collectionID,
	}

	roleBindings, err := client.RoleBindingsAPI.CreateCollectionRoleBindings(&collectionRole)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(roleBindings[0].ID)
	return resourceReadCollectionRole(ctx, d, m)
}

func resourceReadCollectionRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] > Reading Collection Role ")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	roleBindingsID := d.Id()
	roleType := d.Get("role_type").(string)
	if !utils.ValidateCollectionRole(roleType) {
		return diag.FromErr(errors.New("invalid role type"))
	}
	entityID := d.Get("entity_id").(string)
	collectionID := d.Get("collection_id").(string)
	role, err := client.RoleBindingsAPI.ReadRoleBindings(roleBindingsID, roleType, entityID, collectionID, "")
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "read collection role failure"))
	}

	d.SetId(role.ID)
	return nil
}
