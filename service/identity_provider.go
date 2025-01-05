// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"fmt"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
)

type IdentityProviderAPI struct {
	cli *sdk.APIClient
}

func NewIdentityProviderAPI(client *sdk.APIClient) *IdentityProviderAPI {
	return &IdentityProviderAPI{
		cli: client,
	}
}

func (u *IdentityProviderAPI) GetIdentityProviderId(identityProviderName string) (string, error) {
	directoryProviders, _, err := u.cli.IdentityProvidersApi.ListIdentityProviders(context.Background(), &sdk.IdentityProvidersApiListIdentityProvidersOpts{IncludeLocal: optional.NewBool(true)})
	if err != nil {
		return "", err
	}

	for _, directoryProvider := range directoryProviders {
		if directoryProvider.Name == identityProviderName {
			return directoryProvider.Id, nil
		}
	}

	return "", errors.Errorf("can't find identity provider with name '%s'", identityProviderName)
}

func (u *IdentityProviderAPI) GetIdentityProviderTypeById(identityProviderId string) (sdk.IdentityProviderType, error) {
	directoryProviders, _, err := u.cli.IdentityProvidersApi.ListIdentityProviders(context.Background(), &sdk.IdentityProvidersApiListIdentityProvidersOpts{IncludeLocal: optional.NewBool(true)})
	if err != nil {
		return "", err
	}

	for _, directoryProvider := range directoryProviders {
		if directoryProvider.Id == identityProviderId {
			return *directoryProvider.Provider, nil
		}
	}

	return "", errors.Errorf("can't find identity provider with id '%s'", identityProviderId)
}

func (u *IdentityProviderAPI) GetUserDisplayNameTypeById(identityProviderId string, IdentifierInProvider string) (string, error) {
	user, _, err := u.cli.UsersApi.GetUser(context.Background(), identityProviderId, IdentifierInProvider)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", user.FirstName, user.LastName), nil
}

func (u *IdentityProviderAPI) GetGroupDisplayNameTypeById(identityProviderId string, IdentifierInProvider string) (string, error) {
	group, _, err := u.cli.GroupsApi.GetGroup(context.Background(), identityProviderId, IdentifierInProvider)
	if err != nil {
		return "", err
	}

	return group.Name, nil
}
