package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
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
	body := sdk.CollectionsBody{Name: name}
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

// listCollections list collections
func (c *CollectionAPI) ListCollections(name string) (*[]dto.Collection, error) {
	var body sdk.CollectionsApiListCollectionsOpts
	if name != "" {
		body.Name = optional.NewString(name)
	}
	collections, _, err := c.cli.CollectionsApi.ListCollections(context.Background(), &body)
	if err != nil {
		return nil, err
	}
	collectionsDTO, err := dto.ConvertCollectionsToDTO(collections.Content)
	if err != nil {
		return nil, err
	}
	return collectionsDTO, nil
}

// CreateTenantRoleBindings create tenant role bindings
func (c *CollectionAPI) CreateTenantRoleBindings(tenantRole sdk.RoleType, entities *[]sdk.DirectoryEntity) (*[]dto.RoleBindings, error) {

	// get root collection id
	collections, err := c.ListCollections("default")
	if err != nil {
		return nil, err
	}
	if len(*collections) == 0 {
		return nil, errors.New("no root collections found")
	}
	rootCollectionID := (*collections)[0].ID
	subjectType := sdk.COLLECTION_SubjectType
	// create role binding body
	roleBindingBody := sdk.CollectionRolebindingsBody{
		Entities:    *entities,
		RoleType:    &tenantRole,
		SubjectType: &subjectType,
		SubjectID:   rootCollectionID.String(),
	}
	// create role bindings
	roleBindings, _, err := c.cli.CollectionsApi.CreateRoleBinding(context.Background(), roleBindingBody)
	if err != nil {
		return nil, err
	}
	// convert role bindings to dto
	roleBindingsDTO, err := dto.ConvertRoleBindingsToDTO(&roleBindings)

	return roleBindingsDTO, err
}
