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

func LuminateCollectionRoles() *schema.Resource {
	return &schema.Resource{

		Schema: map[string]*schema.Schema{
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Role",
			},
			"collection_id": {
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
		Create: resourceCreateCollectionRoleBinding,
		Read:   resourceReadCollectionRoleBinding,
		Delete: resourceDeleteRoleBinding,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateCollectionRoleBinding(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Role Bindings")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	role := d.Get("role").(string)
	roleType, err := validateCollectionBindingType(role)
	if err != nil {
		return errors.Wrap(err, "validate error:")
	}

	entityID := d.Get("entity_id").(string)
	identityProviderID := d.Get("identity_provider_id").(string)
	identityProviderType, err := client.IdentityProviders.GetIdentityProviderTypeById(identityProviderID)
	if err != nil {
		return errors.Wrap(err, "failed to get identity provider type")
	}

	collectionID := d.Get("collection_id").(string)
	utils.ValidateUuid(collectionID, "collection_id")

	entityType := sdk.USER_EntityType

	entity := sdk.DirectoryEntity{
		IdentifierInProvider: entityID,
		IdentityProviderId:   identityProviderID,
		IdentityProviderType: &identityProviderType,
		Type_:                &entityType,
		DisplayName:          "displayName",
	}

	roleBindings, err := client.CollectionAPI.CreateCollectionRoleBinding(roleType, &entity, collectionID)
	if err != nil {
		return errors.Wrap(err, "failed to create role bindings")
	}

	d.SetId(fmt.Sprintf("%s", (*roleBindings)[0].ID))
	d.Set("entity_id", entityID)
	d.Set("identity_provider_id", identityProviderID)
	d.Set("collection_id", collectionID)

	return nil
}

func resourceReadCollectionRoleBinding(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Role Bindings")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	siteID := d.Get("collection_id").(string)
	roles, err := client.CollectionAPI.ListCollectionRoleBindings(siteID)
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

func validateCollectionBindingType(roleType string) (sdk.CollectionRoleType, error) {
	switch roleType {
	case "PolicyOwner":
		return sdk.POLICY_OWNER_CollectionRoleType, nil
	case "PolicyEntityAssigner":
		return sdk.POLICY_ENTITY_ASSIGNER_CollectionRoleType, nil
	case "ApplicationOwner":
		return sdk.APPLICATION_OWNER_CollectionRoleType, nil
	}
	return "", errors.New(fmt.Sprintf("invalid collection role type: %s", roleType))
}
