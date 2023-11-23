package provider

import (
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			Description:  "The entity id to which this role is assigned.",
			Required:     true,
			ValidateFunc: utils.ValidateString,
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

func resourceDeleteRoleBindings(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	roleType := d.Get("role_type").(string)
	log.Println(fmt.Sprintf("[Info] Delete Role binding for role: %s", roleType))
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	roleID := d.Id()

	err := client.RoleBindingsAPI.DeleteRoleBindings(roleID)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed delete tenant role"))
	}

	return nil
}

func getEntityByRoleBindings(client *service.LuminateService, baseFields dto.RoleBinding) (dto.DirectoryEntity, error) {
	identityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(baseFields.EntityIDPID)
	if err != nil {
		return dto.DirectoryEntity{}, errors.Wrapf(err, "Failed to lookup identity provider type for identity provider id %s", baseFields.EntityIDPID)
	}
	displayName, err := client.IdentityProviders.GetUserDisplayNameTypeById(baseFields.EntityIDPID, baseFields.EntityIDInIDP)
	if err != nil {
		return dto.DirectoryEntity{}, errors.Wrapf(err, "Failed to lookup displayName by IDPID and IDInIDP: %s, %s", baseFields.EntityIDPID, baseFields.EntityIDInIDP)
	}
	entity := dto.DirectoryEntity{
		IdentifierInProvider: baseFields.EntityIDInIDP,
		IdentityProviderId:   baseFields.EntityIDPID,
		EntityType:           baseFields.EntityType,
		IdentityProviderType: dto.ConvertIdentityProviderTypeToString(identityProviderType),
		DisplayName:          displayName,
	}
	return entity, nil
}
