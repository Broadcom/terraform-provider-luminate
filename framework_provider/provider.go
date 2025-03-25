package framework_provider

import (
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkSchema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Ensure MyProductFrameworkProvider implements provider.Provider
var _ provider.Provider = &LuminateFrameworkProvider{}

// MyProductFrameworkProvider implements Terraform Plugin Framework provider
type LuminateFrameworkProvider struct {
	primaryProvider *sdkSchema.Provider
}

func ConvertSchemaAttributes(primarySchema map[string]*sdkSchema.Schema) map[string]schema.Attribute {
	frameworkSchema := make(map[string]schema.Attribute)

	for key, attr := range primarySchema {
		switch attr.Type {
		case sdkSchema.TypeString:
			frameworkSchema[key] = schema.StringAttribute{
				Required: attr.Required,
				Optional: attr.Optional,
			}
		case sdkSchema.TypeInt:
			frameworkSchema[key] = schema.Int64Attribute{
				Required: attr.Required,
				Optional: attr.Optional,
			}
		case sdkSchema.TypeBool:
			frameworkSchema[key] = schema.BoolAttribute{
				Required: attr.Required,
				Optional: attr.Optional,
			}
		case sdkSchema.TypeList:
			frameworkSchema[key] = schema.ListAttribute{
				ElementType: types.StringType,
				Required:    attr.Required,
				Optional:    attr.Optional,
			}
		case sdkSchema.TypeSet:
			frameworkSchema[key] = schema.SetAttribute{
				ElementType: types.StringType,
				Required:    attr.Required,
				Optional:    attr.Optional,
			}
		case sdkSchema.TypeMap:
			frameworkSchema[key] = schema.MapAttribute{
				ElementType: types.StringType,
				Required:    attr.Required,
				Optional:    attr.Optional,
			}
		}
	}

	return frameworkSchema
}

func (provider *LuminateFrameworkProvider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
}

func (provider *LuminateFrameworkProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: ConvertSchemaAttributes(provider.primaryProvider.Schema),
		Blocks:     make(map[string]schema.Block),
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
	return nil
}

func (provider *LuminateFrameworkProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewSiteRegistrationKeyEphemeralResource(),
	}
}

func (provider *LuminateFrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func NewLuminateFrameworkProvider(primaryProvider *sdkSchema.Provider) provider.Provider {
	return &LuminateFrameworkProvider{
		primaryProvider: primaryProvider,
	}
}
