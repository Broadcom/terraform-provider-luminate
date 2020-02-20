package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"fmt"
	"golang.org/x/oauth2/clientcredentials"
)

type LuminateService struct {
	cli *sdk.APIClient

	Sites             *SiteAPI
	Connectors        *ConnectorsAPI
	Applications      *ApplicationAPI
	AccessPolicies    *AccessPolicyAPI
	Users             *UserAPI
	Groups            *GroupAPI
	IdentityProviders *IdentityProviderAPI
	IntegrationAPI    *IntegrationAPI
}

const (
	MaxRequestsPerSecond float64 = 1
	MaxBurst             int     = 1
)

func NewClient(ClientID string, ClientSecret string, Endpoint string) *LuminateService {
	tokenURL := fmt.Sprintf("https://%s/v1/oauth/token", Endpoint)
	basePath := fmt.Sprintf("https://%s/v2", Endpoint)

	cfg := clientcredentials.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		TokenURL:     tokenURL,
		Scopes:       []string{},
	}
	httpClient := cfg.Client(context.Background())

	httpClient.Transport = NewRateLimitTransport(MaxRequestsPerSecond, MaxBurst, httpClient.Transport)

	var lumSvc LuminateService

	lumSvc.cli = sdk.NewAPIClient(&sdk.Configuration{
		UserAgent:  "luminate-terraform-provider",
		HTTPClient: httpClient,
		BasePath:   basePath,
	})

	lumSvc.Sites = NewSiteAPI(lumSvc.cli)
	lumSvc.Connectors = NewConnectorsAPI(lumSvc.cli)
	lumSvc.Applications = NewApplicationAPI(lumSvc.cli)
	lumSvc.AccessPolicies = NewAccessPolicyAPI(lumSvc.cli)
	lumSvc.Users = NewUserAPI(lumSvc.cli)
	lumSvc.Groups = NewGroupAPI(lumSvc.cli)
	lumSvc.IdentityProviders = NewIdentityProviderAPI(lumSvc.cli)
	lumSvc.IntegrationAPI = NewIntegrationAPI(lumSvc.cli, httpClient, basePath)

	return &lumSvc
}
