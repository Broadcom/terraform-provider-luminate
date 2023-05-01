package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"errors"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/antihax/optional"
)

type ConnectorsAPI struct {
	cli *sdk.APIClient
}

func NewConnectorsAPI(client *sdk.APIClient) *ConnectorsAPI {
	return &ConnectorsAPI{
		cli: client,
	}
}

func (api *ConnectorsAPI) GetConnectorByID(connectorID string) (*dto.Connector, error) {
	con, resp, err := api.cli.ConnectorsApi.GetConnector(context.Background(), connectorID)

	if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if con.DeploymentType == nil {
		return nil, nil
	}

	return &dto.Connector{
		ID:        con.Id,
		Name:      con.Name,
		K8SVolume: con.KubernetesPersistentVolumeName,
		Enabled:   con.Enabled,
		Type:      ConnectorTypeFromDeployment(*con.DeploymentType),
		OTP:       con.Otp,
	}, nil
}

func (api *ConnectorsAPI) CreateConnector(connector *dto.Connector, siteID string) (*dto.Connector, error) {
	cType := ConnectorTypeFromString(connector.Type)

	conOptBody := sdk.Connector{
		Name:                           connector.Name,
		DeploymentType:                 &cType,
		Enabled:                        connector.Enabled,
		KubernetesPersistentVolumeName: connector.K8SVolume,
	}

	conOpt := sdk.ConnectorsApiCreateConnectorOpts{
		Body: optional.NewInterface(conOptBody),
	}

	nCon, resp, err := api.cli.ConnectorsApi.CreateConnector(context.Background(), siteID, &conOpt)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		if resp.StatusCode != 201 {
			errMsg := fmt.Sprintf("received bad status code creating connector. Status Code: %d", resp.StatusCode)
			return nil, errors.New(errMsg)
		}
	} else {
		return nil, errors.New("received empty response from the server")
	}

	connector.ID = nCon.Id
	connector.OTP = nCon.Otp
	connector.SiteID = siteID
	return connector, nil
}

func (api *ConnectorsAPI) DeleteConnector(connectorID string) error {
	resp, err := api.cli.ConnectorsApi.DeleteConnector(context.Background(), connectorID)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode != 204 {
			errMsg := fmt.Sprintf("received bad status code deleting connector. Status Code: %d", resp.StatusCode)
			return errors.New(errMsg)
		}
	} else {
		return errors.New("received empty response from the server")
	}

	return nil
}

func (api *ConnectorsAPI) GetConnectorCommand(connectorID string) (string, error) {
	cmd, resp, err := api.cli.ConnectorsApi.GetConnectorCommand(context.Background(), connectorID)
	if err != nil {
		return "", err
	}

	if resp != nil {
		if resp.StatusCode != 200 {
			errMsg := fmt.Sprintf("received bad status code getting connector command. Status Code: %d", resp.StatusCode)
			return "", errors.New(errMsg)
		}
	} else {
		return "", errors.New("received empty response from the server")
	}

	return cmd.DeploymentCommands, nil
}

func ConnectorTypeFromDeployment(ct sdk.DeploymentType) string {
	switch ct {
	case sdk.LINUX_DeploymentType:
		return "linux"
	case sdk.WINDOWS_DeploymentType:
		return "windows"
	case sdk.KUBERNETES_DeploymentType:
		return "kubernetes"
	case sdk.DOCKER_COMPOSE_DeploymentType:
		return "docker-compose"
	}
	return ""
}

func ConnectorTypeFromString(ct string) sdk.DeploymentType {
	switch ct {
	case "linux":
		return sdk.LINUX_DeploymentType
	case "windows":
		return sdk.WINDOWS_DeploymentType
	case "kubernetes":
		return sdk.KUBERNETES_DeploymentType
	case "docker-compose":
		return sdk.DOCKER_COMPOSE_DeploymentType
	}
	return ""
}
