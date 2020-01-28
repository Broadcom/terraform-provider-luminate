package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type IntegrationAPI struct {
	cli        *sdk.APIClient
	httpClient *http.Client
	BasePath   string
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
