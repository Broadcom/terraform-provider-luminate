// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
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
	userPage, _, err := u.cli.UsersApi.SearchUsersByIdP(context.Background(), identityProviderId, &sdk.UsersApiSearchUsersByIdPOpts{Email: optional.NewString(email)})
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
