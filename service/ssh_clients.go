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
	page := int32(0)
	for {
		requestParams := &sdk.GetAllSshClientsOpts{
			Filter: optional.NewString(sshClientName),
			Size:   optional.NewInt32(100),
			Page:   optional.NewInt32(page),
		}
		sshClientPage, _, err := g.cli.SSHClientsApi.GetAllSshClients(context.Background(), requestParams)
		if err != nil {
			return nil, err
		}

		if len(sshClientPage.Content) == 0 || sshClientPage.Last {
			return nil, errors.Errorf("can't find ssh client: '%s'", sshClientName)
		}

		for _, sshClient := range sshClientPage.Content {
			if sshClientName == sshClient.Name {
				return &sshClient, nil
			}
		}

		page = page + 1
	}
}
