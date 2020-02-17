package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
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
	groupPage, _, err := g.cli.GroupsApi.SearchGroupsbyIdp(context.Background(), identityProviderId, &sdk.SearchGroupsbyIdpOpts{Filter: optional.NewString(groupName)})
	if err != nil {
		return "", err
	}

	if len(groupPage.Content) < 1 {
		return "", errors.New("no groups found")
	}

	for _, group := range groupPage.Content {
		if groupName == group.Name {
			return group.Id, nil
		}
	}

	return "", errors.Errorf("can't find group: '%s'", groupName)
}
