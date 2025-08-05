package dto

import (
	"errors"
	"github.com/Broadcom/terraform-provider-luminate/service/utils"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
)

func convertTenantRole(role string) *sdk.TenantRoleType {
	var tenantRoleType sdk.TenantRoleType
	switch role {
	case utils.TenantAdmin:
		tenantRoleType = sdk.TENANT_ADMIN_TenantRoleType
	case utils.TenantViewer:
		tenantRoleType = sdk.TENANT_VIEWER_TenantRoleType
	}

	return &tenantRoleType
}

func convertCollectionRole(role string) *sdk.CollectionRoleType {
	var collectionRoleType sdk.CollectionRoleType
	switch role {
	case utils.ApplicationOwner:
		collectionRoleType = sdk.APPLICATION_OWNER_CollectionRoleType
	case utils.PolicyOwner:
		collectionRoleType = sdk.POLICY_OWNER_CollectionRoleType

	}

	return &collectionRoleType

}

func convertSiteRole(role string) *sdk.SiteRoleType {
	var siteRoleType sdk.SiteRoleType
	switch role {
	case utils.SiteEditor:
		siteRoleType = sdk.SITE_EDITOR_SiteRoleType
	case utils.SiteConnectorDeployer:
		siteRoleType = sdk.SITE_CONNECTOR_DEPLOYER_SiteRoleType
	}

	return &siteRoleType

}

func ConvertTenantRoleBindingsToModel(tenantRole *CreateRoleDTO) (*sdk.CollectionTenantrolebindingsBody, error) {
	if len(tenantRole.Entities) == 0 || tenantRole.Role == "" {
		return &sdk.CollectionTenantrolebindingsBody{}, errors.New("invalid input")
	}
	entities := EntityDTOToEntityModel(tenantRole.Entities)
	role := convertTenantRole(tenantRole.Role)
	return &sdk.CollectionTenantrolebindingsBody{
		Entities: entities,
		RoleType: role,
	}, nil
}

func ConvertCollectionRoleBindingsToModel(collectionRole *CreateCollectionRoleDTO) (*sdk.CollectionCollectionrolebindingsBody, error) {
	if len(collectionRole.Entities) == 0 || collectionRole.Role == "" {
		return &sdk.CollectionCollectionrolebindingsBody{}, errors.New("invalid input")
	}
	entities := EntityDTOToEntityModel(collectionRole.Entities)
	role := convertCollectionRole(collectionRole.Role)
	return &sdk.CollectionCollectionrolebindingsBody{
		Entities:     entities,
		RoleType:     role,
		CollectionId: collectionRole.CollectionID,
	}, nil
}

func ConvertSiteRoleBindingsToModel(siteRole *CreateSiteRoleDTO) (*sdk.CollectionSiterolebindingsBody, error) {
	if len(siteRole.Entities) == 0 || siteRole.Role == "" {
		return &sdk.CollectionSiterolebindingsBody{}, errors.New("invalid input")
	}
	entities := EntityDTOToEntityModel(siteRole.Entities)
	role := convertSiteRole(siteRole.Role)
	return &sdk.CollectionSiterolebindingsBody{
		Entities: entities,
		RoleType: role,
		SiteId:   siteRole.SiteID,
	}, nil
}

func ConvertRolesBindingsToDTO(roleBindings sdk.RoleBindings) []*RoleBinding {
	var dtoRoles []*RoleBinding
	for _, role := range roleBindings.RoleBindings {
		dtoRoles = append(dtoRoles, ConvertRolesBindingToDTO(role))
	}
	return dtoRoles
}

func ConvertRolesBindingToDTO(roleBinding sdk.RoleBinding) *RoleBinding {
	return &RoleBinding{
		ID:            roleBinding.Id,
		EntityIDInIDP: roleBinding.Entity.IdentifierInProvider,
		EntityIDPID:   roleBinding.Entity.IdentityProviderId,
		CollectionID:  roleBinding.Collection.Id,
		ResourceID:    roleBinding.Resource.Id,
	}
}
