package provider

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func LuminateTenantRoleBindings() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Role",
			},
			"entities": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User ID",
						},
						"identity_provider_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Identity Provider ID",
						},
					},
				},
			},
		},
		Create: resourceCreateRoleBinding,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceReadRoleBinding(d *schema.ResourceData, m interface{}) error {
	return nil
}

type EntityTerraform struct {
	ID                 string
	IdentityProviderID string
}

func resourceCreateRoleBinding(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Role Bindings")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	role := d.Get("role").(string)
	roleType, err := validateRoleBindingType(role)
	if err != nil {
		return errors.Wrap(err, "validate error:")
	}

	entitiesTF := d.Get("entity").([]EntityTerraform)
	entities := make([]sdk.DirectoryEntity, len(entitiesTF))
	for _, entity := range entitiesTF {
		identityProviderID := entity.IdentityProviderID
		identityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(identityProviderID)
		if err != nil {
			return errors.Wrap(err, "failed to get identity provider type")
		}
		entityType := sdk.USER_EntityType
		entityID := entity.ID
		entity := sdk.DirectoryEntity{
			IdentifierInProvider: entityID,
			IdentityProviderId:   identityProviderID,
			IdentityProviderType: &identityProviderType,
			Type_:                &entityType,
			DisplayName:          "",
		}
		entities = append(entities, entity)
	}

	roleBindings, err := client.CollectionAPI.CreateTenantRoleBindings(roleType, &entities)
	if err != nil {
		return errors.Wrap(err, "failed to create role bindings")
	}

	d.SetId(fmt.Sprintf("%s", (*roleBindings)[0].ID))
	return nil
}

func resourceUpdateRoleBinding(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDeleteRoleBinding(d *schema.ResourceData, m interface{}) error {
	return nil
}

func validateRoleBindingType(roleType string) (sdk.RoleType, error) {
	switch roleType {
	case "TenantAdmin":
		return sdk.TENANT_ADMIN_RoleType, nil
	case "TenantViewer":
		return sdk.TENANT_VIEWER_RoleType, nil
	case "SiteEditor":
		return sdk.SITE_EDITOR_RoleType, nil
	case "SiteConnectorDeployer":
		return sdk.SITE_CONNECTOR_DEPLOYER_RoleType, nil
	case "ApplicationOwner":
		return sdk.APPLICATION_OWNER_RoleType, nil
	case "PolicyOwner":
		return sdk.POLICY_OWNER_RoleType, nil
	case "PolicyEntityAssigner":
		return sdk.POLICY_ENTITY_ASSIGNER_RoleType, nil

	}
	return "", errors.New("invalid role type")
}

func subjectTypeFromRoleBindingType(roleType sdk.RoleType) (sdk.SubjectType, error) {
	switch roleType {
	case sdk.TENANT_ADMIN_RoleType, sdk.TENANT_VIEWER_RoleType:
		return sdk.COLLECTION_SubjectType, nil
	case sdk.SITE_EDITOR_RoleType, sdk.SITE_CONNECTOR_DEPLOYER_RoleType:
		return sdk.SITE_SubjectType, nil
	case sdk.APPLICATION_OWNER_RoleType:
		return sdk.APP_SubjectType, nil
	case sdk.POLICY_OWNER_RoleType, sdk.POLICY_ENTITY_ASSIGNER_RoleType:
		return sdk.POLICY_SubjectType, nil
	}
	return "", errors.New("invalid subject type")
}
