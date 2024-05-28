package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"fmt"
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

	userPage, resp, err := u.cli.UsersApi.SearchUsersbyIdp(context.Background(), identityProviderId, &sdk.UsersApiSearchUsersbyIdpOpts{Email: optional.NewString(email)})
	if resp != nil && (resp.StatusCode == 403 || resp.StatusCode == 404) {
		return "", nil
	}
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

func (u *UserAPI) DeleteUser(identityProviderId string, UserID string) error {
	resp, err := u.cli.UsersApi.DeleteUser(context.Background(), identityProviderId, UserID)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode != 204 {
			errMsg := fmt.Sprintf("received bad status code deleting site. Status Code: %d", resp.StatusCode)
			return errors.New(errMsg)
		}
	} else {
		return errors.New("received empty response from the server")
	}
	return nil
}
