package dto

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	uuid "github.com/google/uuid"
	"strings"
)

func ConvertCollectionSiteLinkToModel(links *[]CollectionSiteLink) []sdk.CollectionSiteLink {
	var linksModel []sdk.CollectionSiteLink
	for _, link := range *links {
		linksModel = append(linksModel, sdk.CollectionSiteLink{
			SiteId:       link.SiteID,
			CollectionId: link.CollectionID,
		})
	}
	return linksModel
}

func ConvertCollectionSiteLinkToDTO(links *[]sdk.CollectionSiteLink) []CollectionSiteLink {
	var linksDto []CollectionSiteLink
	for _, link := range *links {
		linksDto = append(linksDto, CollectionSiteLink{
			SiteID:       link.SiteId,
			CollectionID: link.CollectionId,
		})
	}
	return linksDto
}

func ConvertCollectionToDTO(collection *sdk.Collection) (*Collection, error) {
	collectionID, err := uuid.Parse(collection.Id)
	if err != nil {
		return nil, err
	}
	parentID, err := uuid.Parse(collection.ParentId)
	if err != nil {
		return nil, err
	}
	return &Collection{
		ID:               collectionID,
		Name:             collection.Name,
		ParentId:         parentID,
		CountResources:   collection.CountResources,
		CountLinkedSites: collection.CountLinkedSites,
		Fqdn:             collection.Fqdn,
	}, err
}

func ConvertCollectionsToDTO(collections []sdk.Collection) (*[]Collection, error) {
	var collectionsDTO []Collection
	for _, collection := range collections {
		collectionDTO, err := ConvertCollectionToDTO(&collection)
		if err != nil {
			return nil, err
		}
		collectionsDTO = append(collectionsDTO, *collectionDTO)
	}
	return &collectionsDTO, nil
}

func ConvertRoleBindingsToDTO(roleBindings *sdk.RoleBindings) (*[]RoleBindings, error) {
	var roleBindingsDTO []RoleBindings
	for _, roleBinding := range roleBindings.RoleBindings {
		roleBindingDTO, err := ConvertRoleBindingToDTO(&roleBinding)
		if err != nil {
			return nil, err
		}
		roleBindingsDTO = append(roleBindingsDTO, *roleBindingDTO)
	}
	return &roleBindingsDTO, nil
}

func ConvertRoleBindingToDTO(roleBinding *sdk.RoleBinding) (*RoleBindings, error) {
	entity := DirectoryEntity{
		IdentifierInProvider: roleBinding.Entity.IdentifierInProvider,
		IdentityProviderId:   roleBinding.Entity.IdentityProviderId,
		EntityType:           FromModelType(*roleBinding.Entity.Type_),
		IdentityProviderType: ConvertIdentityProviderTypeToString(roleBinding.Entity.IdentityProviderType),
		DisplayName:          roleBinding.Entity.DisplayName,
	}
	return &RoleBindings{
		ID:       roleBinding.Id,
		RoleType: ConvertRoleTypeToDTO(*roleBinding.Role.RoleType),
		Entity:   entity,
	}, nil
}

func SubjectTypeFromString(subjectType string) sdk.SubjectType {
	switch strings.ToLower(subjectType) {
	case strings.ToLower("Site"):
		return sdk.SITE_SubjectType
	case strings.ToLower("Policy"):
		return sdk.POLICY_SubjectType
	case strings.ToLower("Application"):
		return sdk.APP_SubjectType
	}
	return ""
}

func ConvertRoleTypeToDTO(roleBindingsType sdk.RoleType) string {
	switch roleBindingsType {
	case sdk.POLICY_ENTITY_ASSIGNER_RoleType:
		return "PolicyEntityAssigner"
	case sdk.POLICY_OWNER_RoleType:
		return "PolicyOwner"
	case sdk.APPLICATION_OWNER_RoleType:
		return "ApplicationOwner"
	case sdk.SITE_CONNECTOR_DEPLOYER_RoleType:
		return "SiteConnectorDeployer"
	case sdk.SITE_EDITOR_RoleType:
		return "SiteEditor"
	case sdk.TENANT_ADMIN_RoleType:
		return "TenantAdmin"
	case sdk.TENANT_VIEWER_RoleType:
		return "TenantViewer"
	}
	return ""
}
