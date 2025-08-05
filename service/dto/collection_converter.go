package dto

import (
	uuid "github.com/google/uuid"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
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
