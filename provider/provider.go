// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LUMINATE_API_ENDPOINT", nil),
			},
			"api_client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LUMINATE_API_CLIENT_ID", nil),
			},
			"api_client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LUMINATE_API_CLIENT_SECRET", nil),
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
		apiEndpoint := d.Get("api_endpoint").(string)
		apiClient := d.Get("api_client_id").(string)
		apiSecret := d.Get("api_client_secret").(string)

		cli := service.NewClient(apiClient, apiSecret, apiEndpoint)
		return cli, nil
	}
}
