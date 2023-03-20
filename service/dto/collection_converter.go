package dto

import sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"

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
