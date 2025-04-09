package framework_provider

import (
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	sdkSchema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Ensure MyProductFrameworkProvider implements provider.Provider
var _ provider.Provider = &LuminateFrameworkProvider{}

// MyProductFrameworkProvider implements Terraform Plugin Framework provider
type LuminateFrameworkProvider struct {
	primaryProvider *sdkSchema.Provider
}

func (provider *LuminateFrameworkProvider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "luminate"
}

func (provider *LuminateFrameworkProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	// Provider schemas must be identical across providers
	// Must be identical to terraform-provider-luminate/provider/provider.go schema
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_endpoint": schema.StringAttribute{
				Optional: true,
			},
			"api_client_id": schema.StringAttribute{
				Optional: true,
			},
			"api_client_secret": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: make(map[string]schema.Block),
	}
}

func (provider *LuminateFrameworkProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	v := provider.primaryProvider.Meta()
	service, ok := v.(*service.LuminateService)
	if !ok {
		response.Diagnostics.AddError("Failed to set metadata", "Failed to load LuminateService from primary provider")
		return
	}

	response.DataSourceData = service
	response.ResourceData = service
	response.EphemeralResourceData = service
}

func (provider *LuminateFrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewLuminateCollectionDataSource(),
	}
}

func (provider *LuminateFrameworkProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return nil
}

func (provider *LuminateFrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWebActivityPolicyResource,
	}
}

func NewLuminateFrameworkProvider(primaryProvider *sdkSchema.Provider) provider.Provider {
	return &LuminateFrameworkProvider{
		primaryProvider: primaryProvider,
	}
}
