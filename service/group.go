// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
	"log"
)

type GroupAPI struct {
	cli *sdk.APIClient
}

func NewGroupAPI(client *sdk.APIClient) *GroupAPI {
	return &GroupAPI{
		cli: client,
	}
}

func (g *GroupAPI) GetGroupId(identityProviderId string, groupName string) (string, error) {
	groupPage, _, err := g.cli.GroupsApi.SearchGroupsbyIdp(context.Background(), identityProviderId, &sdk.GroupsApiSearchGroupsbyIdpOpts{Filter: optional.NewString(groupName)})
	if err != nil {
		return "", err
	}

	if len(groupPage.Content) < 1 {
		return "", errors.New(fmt.Sprintf("no groups found with name: %s", groupName))
	}

	for _, group := range groupPage.Content {
		if groupName == group.Name {
			return group.Id, nil
		}
	}

	return "", errors.Errorf("can't find group: '%s'", groupName)
}

func (g *GroupAPI) AssignUser(groupId string, userId string) error {
	_, err := g.cli.GroupsApi.AssignUserToGroup(context.Background(), groupId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (g *GroupAPI) RemoveUser(groupId string, userId string) error {
	_, err := g.cli.GroupsApi.RemoveUserFromGroup(context.Background(), groupId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (g *GroupAPI) CheckAssignedUser(groupId string, userId string) (bool, error) {
	perPage := int32(100)
	offset := int32(0)

	for {
		userPage, _, err := g.cli.GroupsApi.ListAssignedUsers(context.Background(), utils.LocalIdpId, groupId, &sdk.GroupsApiListAssignedUsersOpts{
			PerPage:    optional.NewFloat64(float64(perPage)),
			PageOffset: optional.NewInterface(offset),
		})
		if err != nil {
			return false, err
		}

		if userPage.TotalElements == 0 {
			return false, nil
		}

		for _, user := range userPage.Content {
			if user.Id == userId {
				// user assigned
				return true, nil
			}
		}

		// FIXME: API "Last" always equals true. remove comment after AC-27729 is done
		//if userPage.Last == true {
		//	return false, nil
		//}

		// next page
		offset = offset + perPage
	}
}

func (g *GroupAPI) GetGroup(groupID string, IDPID string) (*dto.Group, error) {
	group, _, err := g.cli.GroupsApi.GetGroup(context.Background(), IDPID, groupID)
	if err != nil {
		return nil, utils.ParseSwaggerError(err)
	}

	return convertGroupToDTO(group), nil
}

// CreateGroup create Group
func (g *GroupAPI) CreateGroup(idpid string, groupName string) (*dto.Group, error) {
	groupsDto := &sdk.Group{Name: groupName}
	body := sdk.GroupsApiCreateGroupOpts{Body: optional.NewInterface(groupsDto)}
	group, _, err := g.cli.GroupsApi.CreateGroup(context.Background(), idpid, &body)
	if err != nil {
		return nil, utils.ParseSwaggerError(err)
	}

	log.Printf("[DEBUG] - Done Creating Group with name %s", groupName)
	return convertGroupToDTO(group), nil
}

func (g *GroupAPI) DeleteGroup(idpid string, groupID string) error {
	_, err := g.cli.GroupsApi.DeleteGroup(context.Background(), idpid, groupID)
	if err != nil {
		return utils.ParseSwaggerError(err)
	}
	return nil
}

func convertGroupToDTO(group sdk.Group) *dto.Group {
	return &dto.Group{
		Name:               group.Name,
		ID:                 group.Id,
		IdentityProviderId: group.IdentityProviderId,
	}
}
