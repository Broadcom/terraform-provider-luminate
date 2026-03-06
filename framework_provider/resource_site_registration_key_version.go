// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package framework_provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SiteRegistrationKeyVersionResource struct {
	BaseLuminateResource
}

type SiteRegistrationKeyVersionModel struct {
	Version        types.Int64 `tfsdk:"version"`
	VersionChanged types.Bool  `tfsdk:"version_changed"`
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
				Required:    true,
				Description: "Version number for the site registration key used by external secrets. Bump this value when you want to rotate the key and update the secret.",
			},
			"version_changed": schema.BoolAttribute{
				Computed:    true,
				Description: "True when the version value changed compared to the previous run. Use this in the luminate_site_registration_key ephemeral so it only rotates when this is true.",
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
	data.VersionChanged = types.BoolValue(true)
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

// Read is called during refresh. We use it to reset version_changed to false.
// This ensures that version_changed is only true for the single apply where the version was actually changed.
func (r *SiteRegistrationKeyVersionResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data SiteRegistrationKeyVersionModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}
	data.VersionChanged = types.BoolValue(false)
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *SiteRegistrationKeyVersionResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data, state SiteRegistrationKeyVersionModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	data.VersionChanged = types.BoolValue(!data.Version.Equal(state.Version))
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *SiteRegistrationKeyVersionResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {

}
