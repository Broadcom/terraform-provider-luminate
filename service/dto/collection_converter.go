package dto

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	uuid "github.com/google/uuid"
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

// ConvertCollectionsToDTO convert collections model to dto
func ConvertCollectionsToDTO(collections *[]sdk.Collection) (*[]Collection, error) {
	var collectionsDTO []Collection
	for _, collection := range *collections {
		collectionDTO, err := ConvertCollectionToDTO(&collection)
		if err != nil {
			return nil, err
		}
		collectionsDTO = append(collectionsDTO, *collectionDTO)
	}
	return &collectionsDTO, nil
}

func ConvertCollectionToModel(collection *Collection) sdk.Collection {
	return sdk.Collection{
		Id:               collection.ID.String(),
		Name:             collection.Name,
		ParentId:         collection.ParentId.String(),
		CountResources:   collection.CountResources,
		CountLinkedSites: collection.CountLinkedSites,
		Fqdn:             collection.Fqdn,
	}
}
