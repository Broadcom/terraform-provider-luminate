package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
)

type SshClientAPI struct {
	cli *sdk.APIClient
}

func NewSshClientAPI(client *sdk.APIClient) *SshClientAPI {
	return &SshClientAPI{
		cli: client,
	}
}

func (g *SshClientAPI) GetSshClientByName(sshClientName string) (*sdk.SshClient, error) {
	page := float64(0)
	for {
		requestParams := &sdk.SSHClientsApiGetAllSshClientsOpts{
			Filter: optional.NewString(sshClientName),
			Size:   optional.NewFloat64(100),
			Page:   optional.NewFloat64(page),
		}
		sshClientPage, _, err := g.cli.SSHClientsApi.GetAllSshClients(context.Background(), requestParams)
		if err != nil {
			return nil, err
		}

		for _, sshClient := range sshClientPage.Content {
			if sshClientName == sshClient.Name {
				return &sshClient, nil
			}
		}

		if len(sshClientPage.Content) == 0 || sshClientPage.Last {
			return nil, errors.Errorf("can't find ssh client: '%s'", sshClientName)
		}
		page++
	}
}
