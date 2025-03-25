package framework_provider

import (
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	FieldID                           = "id"
	FieldSiteID                       = "site_id"
	FieldRevokeExistingKeyImmediately = "revoke_existing_key_immediately"
	FieldToken                        = "token"
)

type SiteRegistrationKeyEphemeralResource struct {
	client *service.LuminateService
}

// ExampleResourceModel describes the resource data model.
type SiteRegistrationKeyEphemeralResourceModel struct {
	ID                          types.String `tfsdk:"id"`
	SiteID                      types.String `tfsdk:"site_id"`
	Token                       types.String `tfsdk:"token"`
	RevokeExistingKeyImminently types.Bool   `tfsdk:"revoke_existing_key_immediately"`
}

func NewSiteRegistrationKeyEphemeralResource() func() ephemeral.EphemeralResource {
	return func() ephemeral.EphemeralResource {
		return &SiteRegistrationKeyEphemeralResource{}
	}
}

func (r *SiteRegistrationKeyEphemeralResource) Metadata(ctx context.Context, request ephemeral.MetadataRequest, response *ephemeral.MetadataResponse) {
	response.TypeName = "luminate_site_registration_key"
}

func (r *SiteRegistrationKeyEphemeralResource) Schema(ctx context.Context, request ephemeral.SchemaRequest, response *ephemeral.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			FieldID: schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the registration key.",
			},
			FieldSiteID: schema.StringAttribute{
				Required:    true,
				Description: "The site ID we want to associate with this registration key",
			},
			FieldToken: schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The secret key that can be used in order to create Connectors",
			},
			FieldRevokeExistingKeyImmediately: schema.BoolAttribute{
				Required:    true,
				Description: "A field to state if the existing registration key should be revoked immediately or be given a 72 hours expiration time",
			},
		},
	}
}

func (r *SiteRegistrationKeyEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*service.LuminateService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *service.LuminateService, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *SiteRegistrationKeyEphemeralResource) Open(ctx context.Context, request ephemeral.OpenRequest, response *ephemeral.OpenResponse) {
	var data SiteRegistrationKeyEphemeralResourceModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	rotateRequest := dto.SiteRegistrationKeyRotateRequest{
		SiteID:            data.SiteID.ValueString(),
		RevokeImmediately: data.RevokeExistingKeyImminently.ValueBool(),
	}

	generatedKey, err := r.client.SitesRegistrationKeys.RotateRegistrationKey(ctx, rotateRequest)
	if err != nil {
		response.Diagnostics.AddError("Failed to generate registration key", err.Error())
		return
	}

	data.ID = types.StringValue(generatedKey.ID)
	data.Token = types.StringValue(generatedKey.Key)

	response.Diagnostics.Append(response.Result.Set(ctx, &data)...)
}

func (r *SiteRegistrationKeyEphemeralResource) Renew(context.Context, ephemeral.RenewRequest, *ephemeral.RenewResponse) {
	fmt.Printf("")
}
