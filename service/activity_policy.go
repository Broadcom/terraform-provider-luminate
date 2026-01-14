// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"encoding/json"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
	"log"
)

type ActivityPolicyAPI struct {
	cli *sdk.APIClient
}

func NewActivityPolicyAPI(client *sdk.APIClient) *ActivityPolicyAPI {
	return &ActivityPolicyAPI{
		cli: client,
	}
}

func (api *ActivityPolicyAPI) CreateActivityPolicy(activityPolicy *dto.ActivityPolicy) (*dto.ActivityPolicy, error) {
	activityPolicyDto := dto.ToActivityPolicyDto(activityPolicy)
	log.Printf("[DEBUG] - Creating Activity Policy")
	ctx := context.Background()
	for i, entity := range activityPolicyDto.DirectoryEntities {
		if *entity.Type_ == sdk.GROUP_EntityType {
			group, _, err := api.cli.GroupsApi.GetAGroup(ctx, entity.IdentityProviderId, entity.IdentifierInProvider)
			if err != nil {
				var genErr sdk.GenericSwaggerError
				if errors.As(err, &genErr) {
					return nil, errors.Wrapf(err, "failed getting group with body error: %s", string(genErr.Body()))
				}
				return nil, errors.Wrapf(err, "failed getting group")
			}
			activityPolicyDto.DirectoryEntities[i].DisplayName = group.Name
		}
	}
	body := sdk.AccessAndActivityPoliciesApiCreateAnAccessOrActivityPolicyOpts{Body: optional.NewInterface(activityPolicyDto)}
	createdActivityPolicyDtoAsMap, _, err := api.cli.AccessAndActivityPoliciesApi.CreateAnAccessOrActivityPolicy(ctx, &body)
	if err != nil {
		var genErr sdk.GenericSwaggerError
		if errors.As(err, &genErr) {
			return nil, errors.Wrapf(err, "failed creating activity policy with name %s with body error: %s",
				activityPolicy.Name, string(genErr.Body()))
		}
		return nil, errors.Wrapf(err, "failed creating activity policy with name %s", activityPolicy.Name)
	}

	createdActivityPolicyDto, err := api.convertActivityPolicyAsMapToActivityPolicyDto(createdActivityPolicyDtoAsMap)
	if err != nil {
		return nil, err
	}

	return dto.FromActivityPolicyDto(createdActivityPolicyDto), nil
}

func (api *ActivityPolicyAPI) UpdateActivityPolicy(activityPolicy *dto.ActivityPolicy) (*dto.ActivityPolicy, error) {
	activityPolicyDto := dto.ToActivityPolicyDto(activityPolicy)
	body := sdk.AccessAndActivityPoliciesApiUpdateAPolicyOpts{Body: optional.NewInterface(activityPolicyDto)}
	updatedActivityPolicyDtoAsMap, _, err := api.cli.AccessAndActivityPoliciesApi.UpdateAPolicy(context.Background(), activityPolicy.Id, &body)
	if err != nil {
		var genErr sdk.GenericSwaggerError
		if errors.As(err, &genErr) {
			return nil, errors.Wrapf(err, "failed updating activity policy with ID %s with body error: %s",
				activityPolicy.Id, string(genErr.Body()))
		}
		return nil, errors.Wrapf(err, "failed updating activity policy with ID %s", activityPolicy.Id)
	}

	updatedActivityPolicyDto, err := api.convertActivityPolicyAsMapToActivityPolicyDto(updatedActivityPolicyDtoAsMap)
	if err != nil {
		return nil, err
	}

	return dto.FromActivityPolicyDto(updatedActivityPolicyDto), nil
}

func (api *ActivityPolicyAPI) GetActivityPolicy(policyId string) (*dto.ActivityPolicy, error) {
	retrievedActivityPolicyDtoAsMap, resp, err := api.cli.AccessAndActivityPoliciesApi.GetAPolicy(context.Background(), policyId)
	if err != nil {
		if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403 || resp.StatusCode == 500) {
			return nil, nil
		}

		return nil, err
	}

	retrievedActivityPolicyDto, err := api.convertActivityPolicyAsMapToActivityPolicyDto(retrievedActivityPolicyDtoAsMap)
	if err != nil {
		return nil, err
	}

	return dto.FromActivityPolicyDto(retrievedActivityPolicyDto), nil
}

func (api *ActivityPolicyAPI) DeleteActivityPolicy(policyId string) error {
	_, err := api.cli.AccessAndActivityPoliciesApi.DeletePolicy(context.Background(), policyId)
	if err != nil {
		var genErr sdk.GenericSwaggerError
		if errors.As(err, &genErr) {
			return errors.Wrapf(err, "failed deleting activity policy with ID %s with body error: %s",
				policyId, string(genErr.Body()))
		}
		return errors.Wrapf(err, "failed deleting activity policy with ID %s", policyId)
	}

	return nil
}

func (api *ActivityPolicyAPI) convertActivityPolicyAsMapToActivityPolicyDto(activityPolicyDtoAsMap interface{}) (sdk.PolicyActivity, error) {

	// Convert the map of interfaces to string
	jsonString, err := json.Marshal(activityPolicyDtoAsMap)
	if err != nil {
		return sdk.PolicyActivity{}, err
	}

	// Use json unmarshal to convert the string  the struct
	activityPolicyDto := sdk.PolicyActivity{}
	err = json.Unmarshal(jsonString, &activityPolicyDto)
	if err != nil {
		return sdk.PolicyActivity{}, err
	}

	return activityPolicyDto, nil
}
