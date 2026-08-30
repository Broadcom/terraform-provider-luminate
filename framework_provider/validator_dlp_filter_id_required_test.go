// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package framework_provider

import (
	"context"
	"testing"

	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newRuleObjectForValidation(t *testing.T, action, dlpFilterID string, dlpFilterIDNull bool) types.Object {
	t.Helper()

	attrTypes := map[string]attr.Type{
		"action":        types.StringType,
		"dlp_filter_id": types.StringType,
	}
	attrValues := map[string]attr.Value{
		"action": types.StringValue(action),
	}
	if dlpFilterIDNull {
		attrValues["dlp_filter_id"] = types.StringNull()
	} else {
		attrValues["dlp_filter_id"] = types.StringValue(dlpFilterID)
	}

	obj, diags := types.ObjectValue(attrTypes, attrValues)
	if diags.HasError() {
		t.Fatalf("failed to build test object: %v", diags)
	}
	return obj
}

func TestDlpFilterIDRequiredValidator(t *testing.T) {
	tests := []struct {
		name            string
		action          string
		dlpFilterID     string
		dlpFilterIDNull bool
		expectError     bool
	}{
		{name: "CDS with filter id", action: dto.DLPCloudDetectionAction, dlpFilterID: "abc", expectError: false},
		{name: "CDS without filter id", action: dto.DLPCloudDetectionAction, dlpFilterIDNull: true, expectError: true},
		{name: "CDS with empty filter id", action: dto.DLPCloudDetectionAction, dlpFilterID: "", expectError: true},
		{name: "TIS_AND_CDS without filter id", action: dto.DLPAndTISAction, dlpFilterIDNull: true, expectError: true},
		{name: "TIS_AND_CDS with filter id", action: dto.DLPAndTISAction, dlpFilterID: "abc", expectError: false},
		{name: "TIS without filter id", action: dto.TISAction, dlpFilterIDNull: true, expectError: false},
		{name: "BLOCK without filter id", action: dto.BlockAction, dlpFilterIDNull: true, expectError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validator.ObjectRequest{
				ConfigValue: newRuleObjectForValidation(t, tt.action, tt.dlpFilterID, tt.dlpFilterIDNull),
			}
			resp := &validator.ObjectResponse{}

			dlpFilterIDRequiredValidator{}.ValidateObject(context.Background(), req, resp)

			if tt.expectError && !resp.Diagnostics.HasError() {
				t.Errorf("expected error, got none")
			}
			if !tt.expectError && resp.Diagnostics.HasError() {
				t.Errorf("expected no error, got: %v", resp.Diagnostics)
			}
		})
	}
}
