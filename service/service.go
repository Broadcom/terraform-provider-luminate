// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service/roundtripper"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
	"golang.org/x/oauth2/clientcredentials"
	"strings"
)

type LuminateService struct {
	cli              *sdk.APIClient
	TenantBaseDomain string

	Sites                 *SiteAPI
	SitesRegistrationKeys *SitesRegistrationKeysAPI
	Connectors            *ConnectorsAPI
	Applications          *ApplicationAPI
	AccessPolicies        *AccessPolicyAPI
	ActivityPolicies      *ActivityPolicyAPI
	Users                 *UserAPI ``
	Groups                *GroupAPI
	IdentityProviders     *IdentityProviderAPI
	IntegrationAPI        *IntegrationAPI
	SshClientApi          *SshClientAPI
	CollectionAPI         *CollectionAPI
	RoleBindingsAPI       *RoleBindingsAPI
	DNSResiliencyAPI      *DNSResiliencyAPI
	SharedObjectAPI       *SharedObjectAPI
}

const (
	MillsBetweenRetries  int     = 1000
	RetrySleepJitter     int     = 250
	MaxRequestsPerSecond float64 = 5
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

	transport := roundtripper.NewSimpleRateLimitTransport(MaxRequestsPerSecond, httpClient.Transport)
	httpClient.Transport = roundtripper.NewRetryableRateLimitTransport(MillsBetweenRetries, RetrySleepJitter, transport)

	var lumSvc LuminateService

	lumSvc.TenantBaseDomain = strings.ReplaceAll(Endpoint, "api.", "")
	lumSvc.cli = sdk.NewAPIClient(&sdk.Configuration{
		UserAgent:  "luminate-terraform-provider",
		HTTPClient: httpClient,
		BasePath:   basePath,
	})

	lumSvc.Sites = NewSiteAPI(lumSvc.cli)
	lumSvc.SitesRegistrationKeys = NewSitesRegistrationKeysAPI(lumSvc.cli)
	lumSvc.Connectors = NewConnectorsAPI(lumSvc.cli)
	lumSvc.CollectionAPI = NewCollectionAPI(lumSvc.cli)
	lumSvc.Applications = NewApplicationAPI(lumSvc.cli)
	lumSvc.AccessPolicies = NewAccessPolicyAPI(lumSvc.cli)
	lumSvc.ActivityPolicies = NewActivityPolicyAPI(lumSvc.cli)
	lumSvc.Users = NewUserAPI(lumSvc.cli)
	lumSvc.Groups = NewGroupAPI(lumSvc.cli)
	lumSvc.IdentityProviders = NewIdentityProviderAPI(lumSvc.cli)
	lumSvc.IntegrationAPI = NewIntegrationAPI(lumSvc.cli, httpClient, basePath)
	lumSvc.SshClientApi = NewSshClientAPI(lumSvc.cli)
	lumSvc.RoleBindingsAPI = NewRoleBindingsAPI(lumSvc.cli)
	lumSvc.DNSResiliencyAPI = NewDNSResiliencyAPI(lumSvc.cli)
	lumSvc.SharedObjectAPI = NewSharedObjectAPI(lumSvc.cli)

	return &lumSvc
}
