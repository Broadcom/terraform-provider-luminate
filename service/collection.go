package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
	"strings"
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

// ListCollections list collections
func (c *CollectionAPI) ListCollections(name string) (*[]dto.Collection, error) {
	var queryParams sdk.CollectionsApiListCollectionsOpts
	if name != "" {
		queryParams.Name = optional.NewString(name)
	}
	collections, _, err := c.cli.CollectionsApi.ListCollections(context.Background(), &queryParams)
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
func (c *CollectionAPI) CreateTenantRoleBindings(tenantRole sdk.TenantRoleType, entity *sdk.DirectoryEntity) (*[]dto.RoleBindings, error) {

	// create role binding body
	roleBindingBody := sdk.CollectionTenantrolebindingsBody{
		Entities: []sdk.DirectoryEntity{*entity},
		RoleType: &tenantRole,
	}
	// create role bindings
	roleBindings, _, err := c.cli.CollectionsApi.CreateTenantRoleBinding(context.Background(), roleBindingBody)
	if err != nil {
		return nil, err
	}
	// convert role bindings to dto
	roleBindingsDTO, err := dto.ConvertRoleBindingsToDTO(&roleBindings)

	return roleBindingsDTO, err
}

// CreateSiteRoleBinding create site role binding
func (c *CollectionAPI) CreateSiteRoleBinding(tenantRole sdk.SiteRoleType, entity *sdk.DirectoryEntity, siteID string) (*[]dto.RoleBindings, error) {
	// create role binding body
	roleBindingBody := sdk.CollectionSiterolebindingsBody{
		Entities: []sdk.DirectoryEntity{*entity},
		RoleType: &tenantRole,
		SiteId:   siteID,
	}
	// create role bindings
	roleBindings, _, err := c.cli.CollectionsApi.CreateSiteRoleBinding(context.Background(), roleBindingBody)
	if err != nil {
		return nil, err
	}
	// convert role bindings to dto
	roleBindingsDTO, err := dto.ConvertRoleBindingsToDTO(&roleBindings)

	return roleBindingsDTO, err
}

// CreateCollectionRoleBinding create site role binding
func (c *CollectionAPI) CreateCollectionRoleBinding(tenantRole sdk.CollectionRoleType, entity *sdk.DirectoryEntity, siteID string) (*[]dto.RoleBindings, error) {
	// create role binding body
	roleBindingBody := sdk.CollectionCollectionrolebindingsBody{
		Entities:     []sdk.DirectoryEntity{*entity},
		RoleType:     &tenantRole,
		CollectionId: siteID,
	}
	// create role bindings
	roleBindings, _, err := c.cli.CollectionsApi.CreateCollectionRoleBinding(context.Background(), roleBindingBody)
	if err != nil {
		return nil, err
	}
	// convert role bindings to dto
	roleBindingsDTO, err := dto.ConvertRoleBindingsToDTO(&roleBindings)

	return roleBindingsDTO, err
}

// ListRoleBindings list role bindings
func (c *CollectionAPI) ListTenantRoleBindings() (*[]dto.RoleBindings, error) {
	rootCollectionID, err := c.GetRootCollectionID()
	if err != nil {
		return nil, err
	}
	subjectType := sdk.COLLECTION_SubjectType
	queryParams := sdk.CollectionsApiListRoleBindingsOpts{
		SubjectId:   optional.NewInterface(rootCollectionID),
		SubjectType: optional.NewInterface(subjectType),
	}
	res, _, err := c.cli.CollectionsApi.ListRoleBindings(context.Background(), &queryParams)
	if err != nil {
		return nil, err
	}
	roleBindings := sdk.RoleBindings{RoleBindings: res.Content}
	roleBindingsDTO, err := dto.ConvertRoleBindingsToDTO(&roleBindings)

	return roleBindingsDTO, err
}

// ListSiteRoleBindings list site role bindings
func (c *CollectionAPI) ListSiteRoleBindings(siteID string) (*[]dto.RoleBindings, error) {
	subjectType := sdk.SITE_SubjectType
	queryParams := sdk.CollectionsApiListRoleBindingsOpts{
		SubjectId:   optional.NewInterface(siteID),
		SubjectType: optional.NewInterface(subjectType),
	}
	res, _, err := c.cli.CollectionsApi.ListRoleBindings(context.Background(), &queryParams)
	if err != nil {
		return nil, err
	}
	roleBindings := sdk.RoleBindings{RoleBindings: res.Content}
	roleBindingsDTO, err := dto.ConvertRoleBindingsToDTO(&roleBindings)

	return roleBindingsDTO, err
}

// ListCollectionRoleBindings list site role bindings
func (c *CollectionAPI) ListCollectionRoleBindings(collectionID string) (*[]dto.RoleBindings, error) {
	subjectType := sdk.COLLECTION_SubjectType
	queryParams := sdk.CollectionsApiListRoleBindingsOpts{
		SubjectId:   optional.NewInterface(collectionID),
		SubjectType: optional.NewInterface(subjectType),
	}
	res, _, err := c.cli.CollectionsApi.ListRoleBindings(context.Background(), &queryParams)
	if err != nil {
		return nil, err
	}
	roleBindings := sdk.RoleBindings{RoleBindings: res.Content}
	roleBindingsDTO, err := dto.ConvertRoleBindingsToDTO(&roleBindings)

	return roleBindingsDTO, err
}

func (c *CollectionAPI) GetRootCollectionID() (string, error) {
	collections, err := c.ListCollections("default")
	if err != nil {
		return "", err
	}
	if len(*collections) == 0 {
		return "", errors.New("no root collection found")
	}
	rootCollectionID := (*collections)[0].ParentId
	return rootCollectionID.String(), nil
}

// DeleteRoleBinding delete role binding
func (c *CollectionAPI) DeleteRoleBinding(roleBindingID string) error {
	body := sdk.RolebindingsDeleteBody{
		RoleBindingIds: []string{roleBindingID},
	}
	_, err := c.cli.CollectionsApi.DeleteRoleBinding(context.Background(), body)
	if err != nil {
		genericErr := err.(sdk.GenericSwaggerError)
		errModel := genericErr.Model().(sdk.ModelApiResponse)
		if strings.Contains(errModel.Message, "last tenant admin") {
			return errors.Wrapf(err, "cannot delete last tenant admin, id: %s, ", roleBindingID)
		}
		return errors.Wrapf(err, "id: %s", roleBindingID)
	}
	return nil
}
