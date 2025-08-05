package service

import (
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
)

type RoleBindingsAPI struct {
	cli *sdk.APIClient
}

func NewRoleBindingsAPI(client *sdk.APIClient) *RoleBindingsAPI {
	return &RoleBindingsAPI{cli: client}
}

// CreateTenantRoleBindings assign tenant role to admin
func (r *RoleBindingsAPI) CreateTenantRoleBindings(tenantRole *dto.CreateRoleDTO) ([]*dto.RoleBinding, error) {
	body, err := dto.ConvertTenantRoleBindingsToModel(tenantRole)
	if err != nil {
		return nil, err
	}
	roles, _, err := r.cli.CollectionsApi.CreateTenantRoleBinding(context.Background(), *body)
	if err != nil {
		return nil, err
	}
	return dto.ConvertRolesBindingsToDTO(roles), nil
}

// CreateCollectionRoleBindings assign tenant role to admin
func (r *RoleBindingsAPI) CreateCollectionRoleBindings(roleDTO *dto.CreateCollectionRoleDTO) ([]*dto.RoleBinding, error) {
	body, err := dto.ConvertCollectionRoleBindingsToModel(roleDTO)
	if err != nil {
		return nil, err
	}

	roles, _, err := r.cli.CollectionsApi.CreateCollectionRoleBinding(context.Background(), *body)
	if err != nil {
		return nil, err
	}
	return dto.ConvertRolesBindingsToDTO(roles), nil
}

// CreateSiteRoleBindings assign site role to admin
func (r *RoleBindingsAPI) CreateSiteRoleBindings(roleDTO *dto.CreateSiteRoleDTO) ([]*dto.RoleBinding, error) {
	body, err := dto.ConvertSiteRoleBindingsToModel(roleDTO)
	if err != nil {
		return nil, err
	}

	roles, _, err := r.cli.CollectionsApi.CreateSiteRoleBinding(context.Background(), *body)
	if err != nil {
		return nil, err
	}
	return dto.ConvertRolesBindingsToDTO(roles), nil
}

// ReadRoleBindings get role
func (r *RoleBindingsAPI) ReadRoleBindings(
	roleID string,
	roleType string,
	entityId string,
	collectionID string,
	siteID string,
) (*dto.RoleBinding, error) {
	params := sdk.CollectionsApiListRoleBindingsOpts{
		EntityIdInIdp: optional.NewString(entityId),
		RoleType:      optional.NewInterface(roleType),
		SubjectType:   optional.NewInterface("Collection"),
		SubjectId:     optional.NewInterface(collectionID),
		Size:          optional.NewFloat64(100),
	}
	if collectionID == "" {
		params.SubjectId = optional.NewInterface(utils.RootCollection)
	}
	if siteID != "" {
		params.SubjectType = optional.NewInterface("Site")
		params.SubjectId = optional.NewInterface(siteID)
	}

	resp, _, err := r.cli.CollectionsApi.ListRoleBindings(context.Background(), &params)
	if err != nil {
		return &dto.RoleBinding{}, errors.Wrap(err, "failed to list role bindings")
	}

	for _, role := range resp.Content {
		if role.Id == roleID {
			return dto.ConvertRolesBindingToDTO(role), err
		}
	}
	return &dto.RoleBinding{}, errors.New("couldn't find role binding")
}

// DeleteRoleBindings delete role bindings by id
func (r *RoleBindingsAPI) DeleteRoleBindings(roleID string) error {
	body := sdk.RolebindingsDeleteBody{RoleBindingIds: []string{roleID}}

	_, err := r.cli.CollectionsApi.DeleteRoleBinding(context.Background(), body)
	if err != nil {
		return errors.Wrap(err, "failed to delete role bindings")
	}
	return nil
}
