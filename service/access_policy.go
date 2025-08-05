package service

import (
	"context"
	"encoding/json"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/antihax/optional"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
	"log"
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
	log.Printf("[DEBUG] - Creating Policy")
	ctx := context.Background()
	for i, entity := range accessPolicyDto.DirectoryEntities {
		if *entity.Type_ == sdk.GROUP_EntityType {
			group, _, err := api.cli.GroupsApi.GetGroup(ctx, entity.IdentityProviderId, entity.IdentifierInProvider)
			if err != nil {
				return nil, err
			}
			accessPolicyDto.DirectoryEntities[i].DisplayName = group.Name
		}
	}
	body := sdk.AccessAndActivityPoliciesApiCreatePolicyOpts{Body: optional.NewInterface(accessPolicyDto)}
	createdAccessPolicyDtoAsMap, _, err := api.cli.AccessAndActivityPoliciesApi.CreatePolicy(ctx, &body)
	if err != nil {
		return nil, err
	}

	createdAccessPolicyDto, err := api.convertAccessPolicyAsMapToAccessPolicyDto(createdAccessPolicyDtoAsMap)
	if err != nil {
		return nil, err
	}

	return dto.ConvertFromDto(createdAccessPolicyDto), nil
}

func (api *AccessPolicyAPI) UpdateAccessPolicy(accessPolicy *dto.AccessPolicy) (*dto.AccessPolicy, error) {
	accessPolicyDto := dto.ConvertToDto(accessPolicy)
	body := sdk.AccessAndActivityPoliciesApiUpdatePolicyOpts{Body: optional.NewInterface(accessPolicyDto)}
	updatedAccessPolicyDtoAsMap, _, err := api.cli.AccessAndActivityPoliciesApi.UpdatePolicy(context.Background(), accessPolicy.Id, &body)
	if err != nil {
		return nil, err
	}

	updatedAccessPolicyDto, err := api.convertAccessPolicyAsMapToAccessPolicyDto(updatedAccessPolicyDtoAsMap)
	if err != nil {
		return nil, err
	}

	return dto.ConvertFromDto(updatedAccessPolicyDto), nil
}

func (api *AccessPolicyAPI) GetAccessPolicy(policyId string) (*dto.AccessPolicy, error) {
	retrievedAccessPolicyDtoAsMap, resp, err := api.cli.AccessAndActivityPoliciesApi.GetPolicy(context.Background(), policyId)
	if err != nil {
		if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403 || resp.StatusCode == 500) {
			return nil, nil
		}

		return nil, err
	}

	retrievedAccessPolicyDto, err := api.convertAccessPolicyAsMapToAccessPolicyDto(retrievedAccessPolicyDtoAsMap)
	if err != nil {
		return nil, err
	}

	return dto.ConvertFromDto(retrievedAccessPolicyDto), nil
}

func (api *AccessPolicyAPI) DeleteAccessPolicy(policyId string) error {
	_, err := api.cli.AccessAndActivityPoliciesApi.DeletePolicy(context.Background(), policyId)
	if err != nil {
		return err
	}

	return nil
}

func (api *AccessPolicyAPI) convertAccessPolicyAsMapToAccessPolicyDto(accessPolicyDtoAsMap interface{}) (sdk.PolicyAccess, error) {

	// Convert the map of interfaces to string
	jsonString, err := json.Marshal(accessPolicyDtoAsMap)
	if err != nil {
		return sdk.PolicyAccess{}, err
	}

	// Use json unmarshal to convert the string  the struct
	accessPolicyDto := sdk.PolicyAccess{}
	err = json.Unmarshal(jsonString, &accessPolicyDto)
	if err != nil {
		return sdk.PolicyAccess{}, err
	}

	return accessPolicyDto, nil
}
