package framework_provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider resource satisfies framework interfaces.
var _ datasource.DataSource = &LuminateCollectionDataSource{}

func NewLuminateCollectionDataSource() func() datasource.DataSource {
	return func() datasource.DataSource {
		return &LuminateCollectionDataSource{}
	}
}

type LuminateCollectionDataSource struct {
	BaseLuminateDataSource
}

type CollectionDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *LuminateCollectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

func (d *LuminateCollectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Collection data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Collection id",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Collection name",
				Required:            true,
			},
		},
	}
}

func (d *LuminateCollectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CollectionDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collection, err := d.client.CollectionAPI.GetCollectionByName(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read collection, got error: %s", err))
		return
	}

	data.ID = types.StringValue(collection.ID.String())
	data.Name = types.StringValue(collection.Name)
	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
