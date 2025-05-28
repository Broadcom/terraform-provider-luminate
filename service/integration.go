// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type IntegrationAPI struct {
	cli        *sdk.APIClient
	httpClient *http.Client
	BasePath   string
}

type AWSRequestBody struct {
	Provider             string   `json:"provider"`
	HostnameTagName      string   `json:"hostname_tag_name"`
	Name                 string   `json:"name"`
	AwsExternalID        string   `json:"aws_external_id"`
	ID                   string   `json:"id"`
	LuminateAwsAccountID string   `json:"luminate_aws_account_id"`
	Regions              []string `json:"regions"`
	AwsRoleArn           string   `json:"aws_role_arn"`
}

func NewIntegrationAPI(client *sdk.APIClient, httpClient *http.Client, basePath string) *IntegrationAPI {
	return &IntegrationAPI{
		cli:        client,
		httpClient: httpClient,
		BasePath:   basePath,
	}
}

func (u *IntegrationAPI) GetIntegrationId(integrationName string) (string, error) {
	request, _ := http.NewRequest("GET", u.BasePath+"/cloud-integrations/integrations?noHealthCheck=true", nil)

	resp, err := u.httpClient.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var cloudIntegrations []map[string]interface{}

	err = json.Unmarshal(respBody, &cloudIntegrations)
	if err != nil {
		return "", err
	}

	for _, integration := range cloudIntegrations {
		if val, ok := integration["name"]; ok && val == integrationName {
			return integration["id"].(string), nil
		}
	}

	return "", errors.Errorf("can't find aws integration with name '%s'", integrationName)
}

func (u *IntegrationAPI) CreateAWSIntegration(integrationName string) (*dto.AwsIntegration, error) {
	requestBody := []byte(
		fmt.Sprintf(
			`{"provider":"amazon","hostname_tag_name":"Name","name":"%s"}`, integrationName))

	request, _ := http.NewRequest("POST", u.BasePath+"/cloud-integrations/integrations",
		bytes.NewBuffer(requestBody))

	resp, err := u.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("Unable to create new AWS integration for %s", integrationName))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var awsIntegration *dto.AwsIntegration

	err = json.Unmarshal(respBody, &awsIntegration)
	if err != nil {
		return nil, err
	}

	return awsIntegration, nil

}

func (u *IntegrationAPI) UpdateAWSIntegration(awsBody *AWSRequestBody) (*dto.AwsIntegrationBind, error) {
	b, err := json.Marshal(awsBody)
	request, _ := http.NewRequest("PUT", fmt.Sprintf("%s/cloud-integrations/integrations/%s",
		u.BasePath, awsBody.ID), bytes.NewBuffer(b))

	resp, err := u.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unable to update AWS integration for %s", awsBody.Name))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var awsIntegration *dto.AwsIntegrationBind

	err = json.Unmarshal(respBody, &awsIntegration)
	if err != nil {
		return nil, err
	}

	return awsIntegration, nil
}

func (u *IntegrationAPI) ReadAWSIntegration(integrationID string) (*dto.AwsIntegration, error) {

	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/cloud-integrations/integrations/%s",
		u.BasePath, integrationID), nil)

	resp, err := u.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unable to read AWS integration for ID %s", integrationID))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var awsIntegration *dto.AwsIntegration

	err = json.Unmarshal(respBody, &awsIntegration)
	if err != nil {
		return nil, err
	}

	return awsIntegration, nil
}

func (u *IntegrationAPI) ReadAWSIntegrationBind(integrationID string) (*dto.AwsIntegrationBind, error) {

	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/cloud-integrations/integrations/%s",
		u.BasePath, integrationID), nil)

	resp, err := u.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unable to read AWS integration for ID %s", integrationID))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var awsIntegration *dto.AwsIntegrationBind

	err = json.Unmarshal(respBody, &awsIntegration)
	if err != nil {
		return nil, err
	}

	return awsIntegration, nil
}

func (u *IntegrationAPI) DeleteAWSIntegration(integrationID string) error {

	request, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/cloud-integrations/integrations/%s",
		u.BasePath, integrationID),
		nil)

	resp, err := u.httpClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Unable to delete AWS integration for ID %s", integrationID))
	}

	return nil
}
