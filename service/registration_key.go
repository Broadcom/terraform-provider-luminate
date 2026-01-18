// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"

	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/antihax/optional"
	"github.com/pkg/errors"

	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
)

type SitesRegistrationKeysAPI struct {
	cli *sdk.APIClient
}

func NewSitesRegistrationKeysAPI(client *sdk.APIClient) *SitesRegistrationKeysAPI {
	return &SitesRegistrationKeysAPI{
		cli: client,
	}
}

func (api *SitesRegistrationKeysAPI) RotateRegistrationKey(ctx context.Context, rotateRequest dto.SiteRegistrationKeyRotateRequest) (*dto.GeneratedSiteRegistrationKey, error) {
	rotateOptions := sdk.SiteRegistrationKeysApiRotateSiteRegistrationKeysOpts{
		Body: optional.NewInterface(sdk.RotateKeyRequestPostBody{RevokeImmediately: rotateRequest.RevokeImmediately}),
	}

	rotationResponse, httpResponse, err := api.cli.SiteRegistrationKeysApi.RotateSiteRegistrationKeys(ctx, rotateRequest.SiteID, &rotateOptions)
	if err != nil {
		var genErr sdk.GenericSwaggerError
		if errors.As(err, &genErr) {
			return nil, errors.Wrapf(err, "error rotating registration key for site %s with body error: %s", rotateRequest.SiteID, string(genErr.Body()))
		}

		return nil, errors.Wrapf(err, "error rotating registration key for site %s", rotateRequest.SiteID)
	}

	if httpResponse == nil {
		return nil, errors.Errorf("received empty response from the server when generating registration key for site %s", rotateRequest.SiteID)
	}

	if httpResponse.StatusCode != 201 {
		return nil, errors.Errorf("received bad status code (expected 201) when generating registration key for site %s (Status Code: %d)", rotateRequest.SiteID, httpResponse.StatusCode)
	}

	return &dto.GeneratedSiteRegistrationKey{
		ID:  rotationResponse.RegistrationKeyId,
		Key: rotationResponse.RegistrationKey,
	}, nil
}
