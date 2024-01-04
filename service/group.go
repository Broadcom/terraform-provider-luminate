package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
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
