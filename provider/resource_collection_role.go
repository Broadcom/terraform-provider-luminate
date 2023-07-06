package provider

import (
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
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
		Schema: collectionSchema,
		Create: resourceCreateCollectionRole,
		Read:   resourceReadCollectionRole,
		Delete: resourceDeleteRoleBindings,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateCollectionRole(d *schema.ResourceData, m interface{}) error {
	log.Println("[Info] Creating Collection Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client ")
	}

	baseFields, err := extractBaseFields(d)
	if err != nil {
		return errors.Wrap(err, "extract collection role base fields error")
	}
	if !utils.ValidateCollectionRole(baseFields.RoleType) {
		return errors.New("invalid role type")
	}

	collectionID := d.Get("collection_id").(string)

	entity, err := getEntityByRoleBindings(client, baseFields)
	if err != nil {
		return err
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
		return err
	}

	d.SetId(roleBindings[0].ID)
	return nil
}

func resourceReadCollectionRole(d *schema.ResourceData, m interface{}) error {
	log.Println("[Info] > Reading Collection Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	roleBindingsID := d.Id()
	roleType := d.Get("role_type").(string)
	if !utils.ValidateCollectionRole(roleType) {
		return errors.New("invalid role type")
	}
	entityID := d.Get("entity_id").(string)
	collectionID := d.Get("collection_id").(string)
	role, err := client.RoleBindingsAPI.ReadRoleBindings(roleBindingsID, roleType, entityID, collectionID, "")
	if err != nil {
		return errors.Wrap(err, "read collection role failure")
	}

	d.SetId(role.ID)
	return nil
}
