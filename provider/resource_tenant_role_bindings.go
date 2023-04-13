package provider

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func LuminateTenantRoles() *schema.Resource {
	return &schema.Resource{

		Schema: map[string]*schema.Schema{
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Role",
			},
			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User ID",
			},
			"identity_provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity Provider ID",
			},
		},
		Create: resourceCreateRoleBinding,
		Read:   resourceReadRoleBinding,
		Delete: resourceDeleteRoleBinding,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
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

	entityID := d.Get("entity_id").(string)
	identityProviderID := d.Get("identity_provider_id").(string)
	identityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(identityProviderID)
	if err != nil {
		return errors.Wrap(err, "failed to get identity provider type")
	}
	entityType := sdk.USER_EntityType

	entity := sdk.DirectoryEntity{
		IdentifierInProvider: entityID,
		IdentityProviderId:   identityProviderID,
		IdentityProviderType: &identityProviderType,
		Type_:                &entityType,
		DisplayName:          "displayName",
	}

	roleBindings, err := client.CollectionAPI.CreateTenantRoleBindings(roleType, &entity)
	if err != nil {
		return errors.Wrap(err, "failed to create role bindings")
	}

	d.SetId(fmt.Sprintf("%s", (*roleBindings)[0].ID))
	d.Set("entity_id", entityID)
	d.Set("identity_provider_id", identityProviderID)

	return nil
}

func resourceReadRoleBinding(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Role Bindings")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}

	roles, err := client.CollectionAPI.ListRoleBindings("asd", "dads")
	if err != nil {
		return errors.Wrap(err, "failed to get role bindings")
	}
	if len(*roles) == 0 {
		d.SetId("")
		return nil
	}
	// find the role binding that matches id
	for _, bindings := range *roles {
		if bindings.ID == d.Id() {
			d.SetId(bindings.ID)
			break
		}
	}
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
