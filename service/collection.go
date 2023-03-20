package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
)

type CollectionAPI struct {
	cli *sdk.APIClient
}

func NewCollectionAPI(client *sdk.APIClient) *CollectionAPI {
	return &CollectionAPI{
		cli: client,
	}
}

// LinkSiteToCollection link sites to collections
func (c *CollectionAPI) LinkSiteToCollection(links []dto.CollectionSiteLink) (*[]dto.CollectionSiteLink, error) {
	modelLinks := dto.ConvertCollectionSiteLinkToModel(&links)
	body := sdk.MasCollectionsitelinksBody{Links: modelLinks}
	createdLinks, _, err := c.cli.CollectionsApi.CreateCollectionSiteLink(context.Background(), body)
	if err != nil {
		return nil, err
	}
	dtoLinks := dto.ConvertCollectionSiteLinkToDTO(&createdLinks.Links)
	return &dtoLinks, nil
}

// UnlinkSiteFromCollection unlink site from collection
func (c *CollectionAPI) UnlinkSiteFromCollection(link dto.CollectionSiteLink) error {
	_, err := c.cli.CollectionsApi.DeleteCollectionSiteLink(context.Background(), link.CollectionID, link.SiteID, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetCollectionSiteLinks get collection site links
func (c *CollectionAPI) GetCollectionSiteLinks(collectionID string) (*[]dto.CollectionSiteLink, error) {
	links, _, err := c.cli.CollectionsApi.GetCollectionSiteLinks(context.Background(), collectionID)
	if err != nil {
		return nil, err
	}
	dtoLinks := dto.ConvertCollectionSiteLinkToDTO(&links.Links)
	return &dtoLinks, nil
}
