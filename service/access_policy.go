package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"bitbucket.org/accezz-io/terraform-provider-symcsc/service/dto"
	"context"
)

type AccessPolicyAPI struct {
	cli *sdk.APIClient
}

func NewAccessPolicyAPI(client *sdk.APIClient) *AccessPolicyAPI {
	return &AccessPolicyAPI{
		cli: client,
	}
}

func (api *AccessPolicyAPI) CreateAccessPolicy(accessPolicy *dto.AccessPolicy) (*dto.AccessPolicy, error) {
	accessPolicyDto := dto.ConvertToDto(accessPolicy)

	accessPolicyInterface, _, err := api.cli.AccessAndActivityPoliciesPreviewApi.CreatePolicy(context.Background(), accessPolicyDto)
	if err != nil {
		return nil, err
	}

	return dto.ConvertFromDto(accessPolicyInterface.(sdk.PolicyAccess)), nil
}

func (api *AccessPolicyAPI) UpdateAccessPolicy(accessPolicy *dto.AccessPolicy) (*dto.AccessPolicy, error) {
	accessPolicyDto := dto.ConvertToDto(accessPolicy)

	accessPolicyInterface, _, err := api.cli.AccessAndActivityPoliciesPreviewApi.UpdatePolicy(context.Background(), accessPolicy.Id, accessPolicyDto)
	if err != nil {
		return nil, err
	}

	return dto.ConvertFromDto(accessPolicyInterface.(sdk.PolicyAccess)), nil
}

func (api *AccessPolicyAPI) GetAccessPolicy(policyId string) (*dto.AccessPolicy, error) {
	accessPolicyInterface, resp, err := api.cli.AccessAndActivityPoliciesPreviewApi.GetPolicy(context.Background(), policyId)
	if err != nil {
		if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 500) {
			return nil, nil
		}

		return nil, err
	}

	return dto.ConvertFromDto(accessPolicyInterface.(sdk.PolicyAccess)), nil
}

func (api *AccessPolicyAPI) DeleteAccessPolicy(policyId string) error {
	_, err := api.cli.AccessAndActivityPoliciesPreviewApi.DeletePolicy(context.Background(), policyId)
	if err != nil {
		return err
	}

	return nil
}
