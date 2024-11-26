package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
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
	group, resp, err := g.cli.GroupsApi.GetGroup(context.Background(), IDPID, groupID)
	if resp != nil && resp.StatusCode == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, utils.ParseSwaggerError(err)
	}
	groupDto := convertGroupToDTO(group)

	return &groupDto, nil
}

// CreateGroup create Group
func (g *GroupAPI) CreateGroup(idpid string, groupName string) (*dto.Group, error) {
	groupsDto := &sdk.Group{Name: groupName}
	body := sdk.GroupsApiCreateGroupOpts{Body: optional.NewInterface(groupsDto)}
	group, resp, err := g.cli.GroupsApi.CreateGroup(context.Background(), idpid, &body)
	if err != nil {
		return nil, utils.ParseSwaggerError(err)
	}

	if resp != nil {
		if resp.StatusCode != 201 {
			errMsg := fmt.Sprintf("received bad status code for creating group. Status Code: %d, group name %s",
				resp.StatusCode, groupName)
			return nil, errors.New(errMsg)
		}
	} else {
		return nil, errors.New("received empty response from the server for creating group")
	}
	log.Printf("[DEBUG] - Done Creating Group")
	groupDto := convertGroupToDTO(group)
	return &groupDto, nil
}

func (g *GroupAPI) DeleteGroup(idpid string, groupID string) error {
	resp, err := g.cli.GroupsApi.DeleteGroup(context.Background(), idpid, groupID)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode != 204 {
			errMsg := fmt.Sprintf("received bad status code for creating group. Status Code: %d, group ID %s",
				resp.StatusCode, groupID)
			return errors.New(errMsg)
		}
	} else {
		return errors.New("received empty response from the server for deleting group ")
	}
	return nil
}

func convertGroupToDTO(group sdk.Group) dto.Group {
	return dto.Group{
		Name:               group.Name,
		ID:                 group.Id,
		IdentityProviderId: group.IdentityProviderId,
	}
}
