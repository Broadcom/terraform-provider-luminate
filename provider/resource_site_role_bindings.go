package provider

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func LuminateSiteRoles() *schema.Resource {
	return &schema.Resource{

		Schema: map[string]*schema.Schema{
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Role",
			},
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID",
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
		Create: resourceCreateSiteRoleBinding,
		Read:   resourceReadSiteRoleBinding,
		Delete: resourceDeleteRoleBinding,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateSiteRoleBinding(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Role Bindings")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	role := d.Get("role").(string)
	roleType, err := validateTenantSiteBindingType(role)
	if err != nil {
		return errors.Wrap(err, "validate error:")
	}

	entityID := d.Get("entity_id").(string)
	identityProviderID := d.Get("identity_provider_id").(string)
	identityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(identityProviderID)
	if err != nil {
		return errors.Wrap(err, "failed to get identity provider type")
	}

	siteID := d.Get("site_id").(string)
	utils.ValidateUuid(siteID, "site_id")

	entityType := sdk.USER_EntityType

	entity := sdk.DirectoryEntity{
		IdentifierInProvider: entityID,
		IdentityProviderId:   identityProviderID,
		IdentityProviderType: &identityProviderType,
		Type_:                &entityType,
		DisplayName:          "displayName",
	}

	roleBindings, err := client.CollectionAPI.CreateSiteRoleBinding(roleType, &entity, siteID)
	if err != nil {
		return errors.Wrap(err, "failed to create role bindings")
	}

	d.SetId(fmt.Sprintf("%s", (*roleBindings)[0].ID))
	d.Set("entity_id", entityID)
	d.Set("identity_provider_id", identityProviderID)

	return nil
}

func resourceReadSiteRoleBinding(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Role Bindings")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	siteID := d.Get("site_id").(string)
	roles, err := client.CollectionAPI.ListSiteRoleBindings(siteID)
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

func validateTenantSiteBindingType(roleType string) (sdk.SiteRoleType, error) {
	switch roleType {
	case "SiteEditor":
		return sdk.SITE_EDITOR_SiteRoleType, nil
	case "TenantViewer":
		return sdk.SITE_CONNECTOR_DEPLOYER_SiteRoleType, nil
	}
	return "", errors.New(fmt.Sprintf("invalid site role type: %s", roleType))
}
