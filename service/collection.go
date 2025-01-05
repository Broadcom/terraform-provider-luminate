// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

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
	body := sdk.CollectionSitelinksBody{Links: modelLinks}
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

// CreateCollection create collection
func (c *CollectionAPI) CreateCollection(name string) (*dto.Collection, error) {
	body := sdk.CollectionBody{Name: name}
	collection, _, err := c.cli.CollectionsApi.CreateCollection(context.Background(), body)
	if err != nil {
		return nil, err
	}
	collectionDTO, err := dto.ConvertCollectionToDTO(&collection)
	if err != nil {
		return nil, err
	}
	return collectionDTO, err
}

// GetCollection get collection by id
func (c *CollectionAPI) GetCollection(collectionID string) (*dto.Collection, error) {
	collection, _, err := c.cli.CollectionsApi.GetCollection(context.Background(), collectionID)
	if err != nil {
		return nil, err
	}
	collectionDTO, err := dto.ConvertCollectionToDTO(&collection)
	if err != nil {
		return nil, err
	}
	return collectionDTO, err
}

// UpdateCollection update collection
func (c *CollectionAPI) UpdateCollection(name string, collectionID string) (*dto.Collection, error) {
	body := sdk.CollectionCollectionidBody{Name: name}
	collection, _, err := c.cli.CollectionsApi.UpdateCollection(context.Background(), body, collectionID)
	if err != nil {
		return nil, err
	}
	collectionDTO, err := dto.ConvertCollectionToDTO(&collection)
	if err != nil {
		return nil, err
	}
	return collectionDTO, err
}

// DeleteCollection delete collection
func (c *CollectionAPI) DeleteCollection(collectionID string) error {
	body := sdk.CollectionsApiDeleteCollectionOpts{}
	_, err := c.cli.CollectionsApi.DeleteCollection(context.Background(), collectionID, &body)
	if err != nil {
		return err
	}
	return nil
}

// GetCollectionsBySite get collections by site
func (c *CollectionAPI) GetCollectionsBySite(siteID string) (*[]string, error) {
	collections, _, err := c.cli.CollectionsApi.GetCollectionsBySite(context.Background(), siteID)
	if err != nil {
		return nil, err
	}
	return &collections.CollectionIds, err
}
