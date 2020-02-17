package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
)

type UserAPI struct {
	cli *sdk.APIClient
}

func NewUserAPI(client *sdk.APIClient) *UserAPI {
	return &UserAPI{
		cli: client,
	}
}

func (u *UserAPI) GetUserId(identityProviderId string, email string) (string, error) {
	userPage, _, err := u.cli.UsersApi.SearchUsersbyIdp(context.Background(), identityProviderId, &sdk.SearchUsersbyIdpOpts{Email: optional.NewString(email)})
	if err != nil {
		return "", err
	}

	if len(userPage.Content) < 1 {
		return "", errors.New("user not found")
	}

	for _, user := range userPage.Content {
		if user.Email == email {
			return user.Id, nil
		}
	}

	return "", errors.Errorf("can't find user with email: '%s'", email)
}
