package framework_provider

import (
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider resource satisfies framework interfaces.
var _ datasource.DataSource = &LuminateSharedObjectDataSource{}

func NewLuminateSharedObjectDataSource() func() datasource.DataSource {
	return func() datasource.DataSource {
		return &LuminateSharedObjectDataSource{}
	}
}

type LuminateSharedObjectDataSource struct {
	BaseLuminateDataSource
}

type SharedObjectDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

func (d *LuminateSharedObjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shared_object"
}

func (d *LuminateSharedObjectDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Shared object data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Shared object id",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Shared object name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Shared object type",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						dto.IpUuid,
						dto.IsolationProfile,
						dto.OpswatGroups,
					),
				},
			},
		},
	}
}

func (d *LuminateSharedObjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SharedObjectDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sharedObjects, err := d.client.SharedObjectAPI.ListSharedObjects("name,asc", 1, 0, data.Name.ValueString(), data.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read shared object, got error: %s", err))
		return
	}
	if len(sharedObjects) == 0 {
		resp.Diagnostics.AddError("Client Error", "Shared object not found")
		return
	}
	data.ID = types.StringValue(sharedObjects[0].ID)
	data.Name = types.StringValue(sharedObjects[0].Name)
	data.Type = types.StringValue(sharedObjects[0].Type)
	tflog.Trace(ctx, "read a shared object data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
