package provider

import (
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func LuminateAssignRoleBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"role_type": {
			Type:         schema.TypeString,
			Description:  "The tenant role type TenantAdmin / TenantViewer",
			Required:     true,
			ValidateFunc: utils.ValidateString,
			ForceNew:     true,
		},
		"identity_provider_id": {
			Type:         schema.TypeString,
			Description:  "The identity provider id",
			Required:     true,
			ValidateFunc: utils.ValidateString,
			ForceNew:     true,
		},
		"entity_id": {
			Type:         schema.TypeString,
			Description:  "The entity id to which this role assigned.",
			Required:     true,
			ValidateFunc: utils.ValidateUuid,
			ForceNew:     true,
		},
		"entity_type": {
			Type:         schema.TypeString,
			Description:  "The type of entity",
			Required:     true,
			ValidateFunc: utils.ValidateString,
			ForceNew:     true,
		},
	}
}

func extractBaseFields(d *schema.ResourceData) (dto.RoleBinding, error) {
	roleID := d.Id()
	roleType := d.Get("role_type").(string)
	entityID := d.Get("entity_id").(string)
	identityProviderID := d.Get("identity_provider_id").(string)
	entityType := d.Get("entity_type").(string)
	if !utils.ValidateEntityType(entityType) {
		return dto.RoleBinding{}, errors.New("invalid entity type")
	}

	return dto.RoleBinding{
		ID:            roleID,
		EntityIDInIDP: entityID,
		EntityIDPID:   identityProviderID,
		EntityType:    entityType,
		RoleType:      roleType,
	}, nil
}

func resourceDeleteRoleBindings(d *schema.ResourceData, m interface{}) error {
	log.Println("[Info] Delete Tenant Role")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	roleID := d.Id()

	err := client.RoleBindingsAPI.DeleteRoleBindings(roleID)
	if err != nil {
		return errors.Wrap(err, "failed delete tenant role")
	}

	return nil
}
