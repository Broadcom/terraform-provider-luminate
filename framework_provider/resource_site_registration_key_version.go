// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package framework_provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

type SiteRegistrationKeyVersionResource struct {
	BaseLuminateResource
}

type SiteRegistrationKeyVersionModel struct {
	Version types.Int64 `tfsdk:"version"`
}

func NewSiteRegistrationKeyVersionResource() resource.Resource {
	return &SiteRegistrationKeyVersionResource{}
}

func (r *SiteRegistrationKeyVersionResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_site_registration_key_version"
}

func (r *SiteRegistrationKeyVersionResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"version": schema.Int64Attribute{
				Computed:    true,
				Description: "Random based timestamp-generated version number of the site registration key used by external secrets",
			},
		},
	}
}

func (r *SiteRegistrationKeyVersionResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data SiteRegistrationKeyVersionModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	newVersion := time.Now().Unix()
	data.Version = types.Int64Value(newVersion)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *SiteRegistrationKeyVersionResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data SiteRegistrationKeyVersionModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.State.RemoveResource(ctx)
}

func (r *SiteRegistrationKeyVersionResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError("Update Not Supported", "Update is not supported for this resource")
}

func (r *SiteRegistrationKeyVersionResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {

}
