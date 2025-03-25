package service

import (
	"context"

	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/antihax/optional"
	"github.com/pkg/errors"

	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
)

type SitesRegistrationKeysAPI struct {
	cli *sdk.APIClient
}

func NewSitesRegistrationKeysAPI(client *sdk.APIClient) *SitesRegistrationKeysAPI {
	return &SitesRegistrationKeysAPI{
		cli: client,
	}
}

func (api *SitesRegistrationKeysAPI) GetSiteRegistrationKeys(ctx context.Context, siteID string) ([]dto.SiteRegistrationKey, error) {
	registrationKeysResponse, httpResponse, err := api.cli.SiteRegistrationKeysApi.GetSiteRegistrationKeys(ctx, siteID)

	if err != nil {
		return nil, errors.Wrapf(err, "error getting registration keys for site %s", siteID)
	}

	if httpResponse == nil {
		return nil, errors.Errorf("received empty response from the server when getting registration keys for site %s", siteID)
	}

	if httpResponse.StatusCode == 404 || httpResponse.StatusCode == 403 {
		return nil, nil
	}

	dtoKeys := make([]dto.SiteRegistrationKey, 0, len(registrationKeysResponse.RegistrationKeys))
	for _, registrationKey := range registrationKeysResponse.RegistrationKeys {
		dtoKeys = append(dtoKeys, dto.SiteRegistrationKey{ID: registrationKey.Id, ExpirationDate: registrationKey.ExpirationDate})
	}

	return dtoKeys, nil
}

func (api *SitesRegistrationKeysAPI) RotateRegistrationKey(ctx context.Context, rotateRequest dto.SiteRegistrationKeyRotateRequest) (*dto.GeneratedSiteRegistrationKey, error) {
	rotateOptions := sdk.SiteRegistrationKeysApiRotateSiteRegistrationKeyOpts{
		Body: optional.NewInterface(sdk.RotateKeyRequestPostBody{RevokeImmediately: rotateRequest.RevokeImmediately}),
	}

	rotationResponse, httpResponse, err := api.cli.SiteRegistrationKeysApi.RotateSiteRegistrationKey(ctx, rotateRequest.SiteID, &rotateOptions)
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

func (api *SitesRegistrationKeysAPI) DeleteRegistrationKeys(ctx context.Context, siteID string) error {
	httpResponse, err := api.cli.SiteRegistrationKeysApi.DeleteSiteRegistrationKeys(ctx, siteID)

	if err != nil {
		return errors.Wrapf(err, "error deleting registration keys for site %s", siteID)
	}

	if httpResponse == nil {
		return errors.Errorf("received empty response from the server when deleting registration keys for site %s", siteID)
	}

	if httpResponse.StatusCode != 204 {
		return errors.Errorf("received bad status code (expected 204) when deleting registration keys for site %s (Status Code: %d)", siteID, httpResponse.StatusCode)
	}

	return nil
}
