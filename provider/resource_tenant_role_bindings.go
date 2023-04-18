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
		Read:   resourceReadTenantRoleBinding,
		Delete: resourceDeleteRoleBinding,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateRoleBinding(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Role Bindings")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	role := d.Get("role").(string)
	roleType, err := validateTenantRoleBindingType(role)
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

func resourceReadTenantRoleBinding(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Role Bindings")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}

	roles, err := client.CollectionAPI.ListTenantRoleBindings()
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
	log.Printf("[INFO] Deleting Role Bindings")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	err := client.CollectionAPI.DeleteRoleBinding(d.Id())
	if err != nil {
		return errors.Wrap(err, "failed to delete role binding")
	}
	return nil
}

func validateTenantRoleBindingType(roleType string) (sdk.TenantRoleType, error) {
	switch roleType {
	case "TenantAdmin":
		return sdk.TENANT_ADMIN_TenantRoleType, nil
	case "TenantViewer":
		return sdk.TENANT_VIEWER_TenantRoleType, nil
	}
	return "", errors.New("invalid tenant role type")
}
