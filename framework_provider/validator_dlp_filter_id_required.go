// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package framework_provider

import (
	"context"
	"fmt"

	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.Object = dlpFilterIDRequiredValidator{}

// dlpFilterIDRequiredValidator ensures dlp_filter_id is set on rules whose action requires it.
type dlpFilterIDRequiredValidator struct{}

func (v dlpFilterIDRequiredValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("dlp_filter_id must be set when action is %q or %q", dto.DLPCloudDetectionAction, dto.DLPAndTISAction)
}

func (v dlpFilterIDRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v dlpFilterIDRequiredValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	action, ok := req.ConfigValue.Attributes()["action"].(types.String)
	if !ok || action.IsNull() || action.IsUnknown() {
		return
	}
	if action.ValueString() != dto.DLPCloudDetectionAction && action.ValueString() != dto.DLPAndTISAction {
		return
	}

	dlpFilterID, ok := req.ConfigValue.Attributes()["dlp_filter_id"].(types.String)
	if !ok || dlpFilterID.IsUnknown() {
		return
	}
	if dlpFilterID.IsNull() || dlpFilterID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing Required Attribute",
			fmt.Sprintf("dlp_filter_id must be set when action is %q or %q.", dto.DLPCloudDetectionAction, dto.DLPAndTISAction),
		)
	}
}
