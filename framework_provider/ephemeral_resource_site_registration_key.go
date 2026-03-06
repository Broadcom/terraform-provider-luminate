// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package framework_provider

import (
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SiteRegistrationKeyEphemeralResource struct {
	BaseLuminateEphemeralResource
}

type SiteRegistrationKeyEphemeralResourceModel struct {
	ID                           types.String `tfsdk:"id"`
	SiteID                       types.String `tfsdk:"site_id"`
	RevokeExistingKeyImmediately types.Bool   `tfsdk:"revoke_existing_key_immediately"`
	IsApplying                   types.Bool   `tfsdk:"is_applying"`
	ShouldRotate                 types.Bool   `tfsdk:"should_rotate"`
	Token                        types.String `tfsdk:"token"`
}

func NewSiteRegistrationKeyEphemeralResource() ephemeral.EphemeralResource {
	return &SiteRegistrationKeyEphemeralResource{}
}

func (r *SiteRegistrationKeyEphemeralResource) Metadata(ctx context.Context, request ephemeral.MetadataRequest, response *ephemeral.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_site_registration_key"
}

func (r *SiteRegistrationKeyEphemeralResource) Schema(ctx context.Context, request ephemeral.SchemaRequest, response *ephemeral.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the registration key.",
			},
			"site_id": schema.StringAttribute{
				Required:    true,
				Description: "The site ID we want to associate with this registration key",
			},
			"revoke_existing_key_immediately": schema.BoolAttribute{
				Required:    true,
				Description: "A field to state if the existing registration key should be revoked immediately or be given a 72 hours expiration time (true: → All existing keys are deleted.)",
			},
			"is_applying": schema.BoolAttribute{
				Required:    true,
				Description: "Set to terraform.applying so the key is only rotated during apply, not during plan. This avoids rotating the key on plan-only runs.",
			},
			"should_rotate": schema.BoolAttribute{
				Required:    true,
				Description: "Set to true when the key should be rotated (e.g. reference version_changed from a luminate_site_registration_key_version resource). The key is only rotated when this is true (and is_applying is true).",
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The secret key that can be used in order to create Connectors",
			},
		},
	}
}

func (r *SiteRegistrationKeyEphemeralResource) Open(ctx context.Context, request ephemeral.OpenRequest, response *ephemeral.OpenResponse) {
	var data SiteRegistrationKeyEphemeralResourceModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Only rotate when both are known and true (avoid ValueBool on unknown at plan)
	shouldRotate := data.IsApplying.ValueBool() && data.ShouldRotate.ValueBool()
	if !shouldRotate {
		data.ID = types.StringNull()
		data.Token = types.StringValue("")
		response.Diagnostics.Append(response.Result.Set(ctx, &data)...)
		return
	}

	rotateRequest := dto.SiteRegistrationKeyRotateRequest{
		SiteID:            data.SiteID.ValueString(),
		RevokeImmediately: data.RevokeExistingKeyImmediately.ValueBool(),
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
