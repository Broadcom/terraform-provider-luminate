// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_client_secret": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"luminate_site":                  LuminateSite(),
			"luminate_connector":             LuminateConnector(),
			"luminate_web_application":       LuminateWebApplication(),
			"luminate_ssh_application":       LuminateSSHApplication(),
			"luminate_ssh_gw_application":    LuminateSshGwApplication(),
			"luminate_tcp_application":       LuminateTCPApplication(),
			"luminate_rdp_application":       LuminateRDPApplication(),
			"luminate_segment_application":   LuminateSegmentApplication(),
			"luminate_web_access_policy":     LuminateWebAccessPolicy(),
			"luminate_tcp_access_policy":     LuminateTcpAccessPolicy(),
			"luminate_ssh_access_policy":     LuminateSshAccessPolicy(),
			"luminate_rdp_access_policy":     LuminateRdpAccessPolicy(),
			"luminate_group_user":            LuminateGroupUser(),
			"luminate_aws_integration":       LuminateAWSIntegration(),
			"luminate_aws_integration_bind":  LuminateAWSIntegrationBind(),
			"luminate_collection_site_link":  LuminateCollectionSiteLink(),
			"luminate_collection":            LuminateCollection(),
			"luminate_tenant_role":           LuminateTenantRole(),
			"luminate_collection_role":       LuminateCollectionRole(),
			"luminate_site_role":             LuminateSiteRole(),
			"luminate_dns_group_resiliency":  LuminateDNSGroupResiliency(),
			"luminate_dns_server_resiliency": LuminateDNSServerResiliency(),
			"luminate_resources_group":       LuminateGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"luminate_group":             LuminateDataSourceGroups(),
			"luminate_user":              LuminateDataSourceUsers(),
			"luminate_identity_provider": LuminateDataSourceIdentityProvider(),
			"luminate_aws_integration":   LuminateDataSourceAwsIntegration(),
			"luminate_ssh_client":        LuminateDataSourceSshClient(),
		},
	}
	p.ConfigureContextFunc = configure()
	return p
}

func configure() func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Due to the migration with terraform-plugin-framework (using the mux package)
		// We need make sure the providers' schema is identical
		// terraform-plugin-framework doesn't allow using "Default" values for "Required" fields
		// Therefore, we set the fields as optional and enforce the "Required" under the configure method

		apiEndpoint := getProviderField(d, "api_endpoint", "LUMINATE_API_ENDPOINT")
		apiClient := getProviderField(d, "api_client_id", "LUMINATE_API_CLIENT_ID")
		apiSecret := getProviderField(d, "api_client_secret", "LUMINATE_API_CLIENT_SECRET")

		if apiEndpoint == "" || apiClient == "" || apiSecret == "" {
			return nil, diag.Errorf("API endpoint, client id, secret are required")
		}

		cli := service.NewClient(apiClient, apiSecret, apiEndpoint)
		return cli, nil
	}
}

func getProviderField(d *schema.ResourceData, key, envName string) string {
	resourceValue := d.Get(key).(string)
	if resourceValue != "" {
		return resourceValue
	}

	return os.Getenv(envName)
}
