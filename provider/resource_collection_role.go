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
		Description:  "the collection id to which this role assigned to.",
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
		return errors.Wrap(err, "extract tenant role base fields error")
	}
	if !utils.ValidateCollectionRole(baseFields.RoleType) {
		return errors.New("invalid role type")
	}

	collectionID := d.Get("collection_id").(string)

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

	var tenantRole = dto.CreateCollectionRoleDTO{
		CreateRoleDTO: dto.CreateRoleDTO{
			Role:     baseFields.RoleType,
			Entities: []dto.DirectoryEntity{entity},
		},
		CollectionID: collectionID,
	}

	roleBindings, err := client.RoleBindingsAPI.CreateCollectionRoleBindings(&tenantRole)
	if err != nil {
		return err
	}

	d.SetId(roleBindings[0].ID)
	return nil
}

func resourceReadCollectionRole(d *schema.ResourceData, m interface{}) error {
	log.Println("[Info] Reading Tenant Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	roleBindingsID := d.Id()
	roleType := d.Get("role_type").(string)
	entityID := d.Get("entity_id").(string)
	collectionID := d.Get("collection_id").(string)
	role, err := client.RoleBindingsAPI.ReadRoleBindings(roleBindingsID, roleType, entityID, collectionID, "")
	if err != nil {
		return errors.Wrap(err, "read tenant role failure")
	}

	d.SetId(role.ID)
	return nil
}
