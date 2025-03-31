package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"encoding/json"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/antihax/optional"
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
			group, _, err := api.cli.GroupsApi.GetGroup(ctx, entity.IdentityProviderId, entity.IdentifierInProvider)
			if err != nil {
				return nil, err
			}
			activityPolicyDto.DirectoryEntities[i].DisplayName = group.Name
		}
	}
	body := sdk.AccessAndActivityPoliciesApiCreatePolicyOpts{Body: optional.NewInterface(activityPolicyDto)}
	createdActivityPolicyDtoAsMap, _, err := api.cli.AccessAndActivityPoliciesApi.CreatePolicy(ctx, &body)
	if err != nil {
		return nil, err
	}

	createdActivityPolicyDto, err := api.convertActivityPolicyAsMapToActivityPolicyDto(createdActivityPolicyDtoAsMap)
	if err != nil {
		return nil, err
	}

	return dto.FromActivityPolicyDto(createdActivityPolicyDto), nil
}

func (api *ActivityPolicyAPI) UpdateActivityPolicy(activityPolicy *dto.ActivityPolicy) (*dto.ActivityPolicy, error) {
	activityPolicyDto := dto.ToActivityPolicyDto(activityPolicy)
	body := sdk.AccessAndActivityPoliciesApiUpdatePolicyOpts{Body: optional.NewInterface(activityPolicyDto)}
	updatedActivityPolicyDtoAsMap, _, err := api.cli.AccessAndActivityPoliciesApi.UpdatePolicy(context.Background(), activityPolicy.Id, &body)
	if err != nil {
		return nil, err
	}

	updatedActivityPolicyDto, err := api.convertActivityPolicyAsMapToActivityPolicyDto(updatedActivityPolicyDtoAsMap)
	if err != nil {
		return nil, err
	}

	return dto.FromActivityPolicyDto(updatedActivityPolicyDto), nil
}

func (api *ActivityPolicyAPI) GetActivityPolicy(policyId string) (*dto.ActivityPolicy, error) {
	retrievedActivityPolicyDtoAsMap, resp, err := api.cli.AccessAndActivityPoliciesApi.GetPolicy(context.Background(), policyId)
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
		return err
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
