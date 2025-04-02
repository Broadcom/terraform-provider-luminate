package framework_provider

import (
	swagger "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BasePolicyResource struct {
	BaseLuminateResource
}

type BasePolicyResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	Name               types.String `tfsdk:"name"`
	IdentityProviderID types.String `tfsdk:"identity_provider_id"`
	UserIDs            types.List   `tfsdk:"user_ids"`
	GroupIDs           types.List   `tfsdk:"group_ids"`
	Applications       types.Set    `tfsdk:"applications"`
	CollectionID       types.String `tfsdk:"collection_id"`
	Conditions         types.Object `tfsdk:"conditions"`
}

func CreatePolicyBaseSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The activity policy ID",
			Computed:    true,
		},
		"enabled": schema.BoolAttribute{
			Description: "Indicates whether this policy is enabled.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(true),
		},
		"name": schema.StringAttribute{
			Description: "A descriptive name of the policy.",
			Required:    true,
		},
		"identity_provider_id": schema.StringAttribute{
			Description: "The identity provider id",
			Optional:    true,
		},
		"user_ids": schema.ListAttribute{
			Description: "The user entities to which this policy applies.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"group_ids": schema.ListAttribute{
			Description: "The group entities to which this policy applies.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"applications": schema.SetAttribute{
			Description: "The applications to which this policy applies.",
			Required:    true,
			ElementType: types.StringType,
		},
		"collection_id": schema.StringAttribute{
			Description: "Collection ID to which the policy will be assigned",
			Computed:    true,
		},
		"conditions": schema.SingleNestedAttribute{
			Description: "Conditions that specify the context of the user and the device in which the policy will apply",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"location": schema.ListAttribute{
					Description: "location based condition, specify the list of accepted locations.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"source_ip": schema.ListAttribute{
					Description: "source ip based condition, specify the allowed CIDR for this policy.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"managed_device": schema.SingleNestedAttribute{
					Description: "list of managed devices that have restriction access",
					Optional:    true,
					Attributes: map[string]schema.Attribute{
						"opswat": schema.BoolAttribute{
							Description: "Indicate whatever to restrict access to Opswat MetaAccess",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"symantec_cloudsoc": schema.BoolAttribute{
							Description: "Indicate whatever to restrict access to symantec cloudsoc",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"symantec_web_security_service": schema.BoolAttribute{
							Description: "Indicate whatever to restrict access to symantec web security service",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
				},
				"unmanaged_device": schema.SingleNestedAttribute{
					Description: "list of unmanaged devices that have restriction access",
					Optional:    true,
					Attributes: map[string]schema.Attribute{
						"opswat": schema.BoolAttribute{
							Description: "Indicate whatever to restrict access to Opswat MetaAccess",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"symantec_cloudsoc": schema.BoolAttribute{
							Description: "Indicate whatever to restrict access to symantec cloudsoc",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"symantec_web_security_service": schema.BoolAttribute{
							Description: "Indicate whatever to restrict access to symantec web security service",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
}

func convertPolicyToBaseModel(ctx context.Context, policy *dto.Policy) (*BasePolicyResourceModel, diag.Diagnostics) {
	policyResourceModel := &BasePolicyResourceModel{
		ID:           types.StringValue(policy.Id),
		Enabled:      types.BoolValue(policy.Enabled),
		Name:         types.StringValue(policy.Name),
		CollectionID: types.StringValue(policy.CollectionID),
		GroupIDs:     types.ListNull(types.StringType),
	}

	appIDs := convertToStringTypeSlice(policy.Applications)
	applications, appDiagnostics := types.SetValueFrom(ctx, types.StringType, appIDs)
	if appDiagnostics.HasError() {
		return nil, appDiagnostics
	}
	policyResourceModel.Applications = applications

	conditions, conditionDiags := flattenConditions(ctx, policy.Conditions)
	if conditionDiags != nil && conditionDiags.HasError() {
		return nil, conditionDiags
	}
	policyResourceModel.Conditions = conditions

	policyResourceModel.GroupIDs = types.ListNull(types.StringType)
	policyResourceModel.UserIDs = types.ListNull(types.StringType)

	if len(policy.DirectoryEntities) > 0 {
		groupIDs := make([]types.String, 0)
		userIDs := make([]types.String, 0)
		for _, entity := range policy.DirectoryEntities {
			if *dto.ToModelType(entity.EntityType) == swagger.GROUP_EntityType {
				groupIDs = append(groupIDs, types.StringValue(entity.IdentifierInProvider))
			}
			if *dto.ToModelType(entity.EntityType) == swagger.USER_EntityType {
				userIDs = append(userIDs, types.StringValue(entity.IdentifierInProvider))
			}
		}

		if policy.DirectoryEntities[0].IdentityProviderId != "" {
			policyResourceModel.IdentityProviderID = types.StringValue(policy.DirectoryEntities[0].IdentityProviderId)
		}

		if len(groupIDs) > 0 {
			groupsIDsList, groupIDsDiagnostics := types.ListValueFrom(ctx, types.StringType, groupIDs)
			if groupIDsDiagnostics.HasError() {
				return nil, groupIDsDiagnostics
			}
			policyResourceModel.GroupIDs = groupsIDsList
		}

		if len(userIDs) > 0 {
			userIDsList, userIDsDiagnostics := types.ListValueFrom(ctx, types.StringType, userIDs)
			if userIDsDiagnostics.HasError() {
				return nil, userIDsDiagnostics
			}
			policyResourceModel.UserIDs = userIDsList
		}
	}
	return policyResourceModel, nil
}

func convertBaseModelToPolicy(ctx context.Context, policyModel *BasePolicyResourceModel) (*dto.Policy, diag.Diagnostics) {
	var directoryEntity []dto.DirectoryEntity
	policy := &dto.Policy{
		Id:           policyModel.ID.ValueString(),
		Name:         policyModel.Name.ValueString(),
		Enabled:      policyModel.Enabled.ValueBool(),
		CollectionID: policyModel.CollectionID.ValueString(),
	}
	if len(policyModel.Applications.Elements()) > 0 {
		diags := policyModel.Applications.ElementsAs(ctx, &policy.Applications, true)
		if diags.HasError() {
			return nil, diags
		}
	}
	if len(policyModel.UserIDs.Elements()) > 0 {
		var userIDs []string
		diags := policyModel.UserIDs.ElementsAs(ctx, &userIDs, true)
		if diags.HasError() {
			return nil, diags
		}
		for _, userID := range userIDs {
			directoryEntity = append(directoryEntity, dto.DirectoryEntity{
				IdentityProviderId:   policyModel.IdentityProviderID.ValueString(),
				IdentifierInProvider: userID,
				EntityType:           "User",
			})
		}
	}

	if len(policyModel.GroupIDs.Elements()) > 0 {
		var groupIDs []string
		diags := policyModel.GroupIDs.ElementsAs(ctx, &groupIDs, true)
		if diags.HasError() {
			return nil, diags
		}
		for _, groupID := range groupIDs {
			directoryEntity = append(directoryEntity, dto.DirectoryEntity{
				IdentityProviderId:   policyModel.IdentityProviderID.ValueString(),
				IdentifierInProvider: groupID,
				EntityType:           "Group",
			})
		}
	}
	policy.DirectoryEntities = directoryEntity

	modelConditions, conditionDiags := extractModelConditions(ctx, policyModel.Conditions)
	if conditionDiags != nil && conditionDiags.HasError() {
		return nil, conditionDiags
	}
	policy.Conditions = modelConditions

	return policy, nil
}

func flattenConditions(ctx context.Context, conditions *dto.Conditions) (types.Object, diag.Diagnostics) {
	emptyConditions := types.ObjectNull(getConditionAttributesTypes())
	if conditions == nil {
		return emptyConditions, nil
	}
	conditionAttributes := map[string]attr.Value{}

	containsCondition := false

	// Location
	conditionAttributes["location"] = types.ListNull(types.StringType)
	if len(conditions.Location) > 0 {
		locationList, diags := types.ListValueFrom(ctx, types.StringType, conditions.Location)
		if diags.HasError() {
			return emptyConditions, diags
		}
		conditionAttributes["location"] = locationList
		containsCondition = true
	}

	// Source IP
	conditionAttributes["source_ip"] = types.ListNull(types.StringType)
	if len(conditions.SourceIp) > 0 {
		sourceIpList, _ := types.ListValueFrom(ctx, types.StringType, conditions.SourceIp)
		conditionAttributes["source_ip"] = sourceIpList
		containsCondition = true
	}

	// Managed Device
	conditionAttributes["managed_device"] = types.ObjectNull(managedDeviceAttributeTypes())
	if hasDeviceCondition(conditions.ManagedDevice) {
		managedDevice := flattenManagedDevice(conditions.ManagedDevice)
		conditionAttributes["managed_device"] = managedDevice
		containsCondition = true
	}

	// Unmanaged Device
	conditionAttributes["unmanaged_device"] = types.ObjectNull(managedDeviceAttributeTypes())
	if hasDeviceCondition(conditions.UnmanagedDevice) {
		unmanagedDevice := flattenManagedDevice(conditions.UnmanagedDevice)
		conditionAttributes["unmanaged_device"] = unmanagedDevice
		containsCondition = true
	}

	if !containsCondition {
		return emptyConditions, nil
	}

	return types.ObjectValue(getConditionAttributesTypes(), conditionAttributes)
}

func getConditionAttributesTypes() map[string]attr.Type {
	conditionTypes := map[string]attr.Type{
		"location":         types.ListType{ElemType: types.StringType},
		"source_ip":        types.ListType{ElemType: types.StringType},
		"managed_device":   types.ObjectType{AttrTypes: managedDeviceAttributeTypes()},
		"unmanaged_device": types.ObjectType{AttrTypes: managedDeviceAttributeTypes()},
	}
	return conditionTypes
}

func convertToStringTypeSlice(slice []string) []types.String {
	values := make([]types.String, 0)
	for _, appID := range slice {
		values = append(values, types.StringValue(appID))
	}
	return values
}

func flattenManagedDevice(manageDevice dto.Device) types.Object {
	deviceAttributes := map[string]attr.Value{}

	deviceAttributes["opswat"] = types.BoolValue(manageDevice.OpswatMetaAccess)
	deviceAttributes["symantec_cloudsoc"] = types.BoolValue(manageDevice.SymantecCloudSoc)
	deviceAttributes["symantec_web_security_service"] = types.BoolValue(manageDevice.SymantecWebSecurityService)

	return types.ObjectValueMust(managedDeviceAttributeTypes(), deviceAttributes)
}

func hasDeviceCondition(managedDevice dto.Device) bool {
	if managedDevice.OpswatMetaAccess {
		return true
	}

	if managedDevice.SymantecCloudSoc {
		return true
	}

	if managedDevice.SymantecWebSecurityService {
		return true
	}
	return false
}

func managedDeviceAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"opswat":                        types.BoolType,
		"symantec_cloudsoc":             types.BoolType,
		"symantec_web_security_service": types.BoolType,
	}
}

func extractModelConditions(ctx context.Context, conditions types.Object) (*dto.Conditions, diag.Diagnostics) {
	if conditions.IsNull() || conditions.IsUnknown() {
		return nil, nil
	}

	result := &dto.Conditions{}

	// Extract Location
	locationValue, ok := conditions.Attributes()["location"]
	if ok && !locationValue.IsNull() && !locationValue.IsUnknown() {
		var locations []string
		diags := locationValue.(types.List).ElementsAs(ctx, &locations, false)
		if diags.HasError() {
			return nil, diags
		}
		result.Location = locations
	}

	// Extract Source IP
	sourceIpValue, ok := conditions.Attributes()["source_ip"]
	if ok && !sourceIpValue.IsNull() && !sourceIpValue.IsUnknown() {
		var sourceIps []string
		diags := sourceIpValue.(types.List).ElementsAs(ctx, &sourceIps, false)
		if diags.HasError() {
			return nil, diags
		}
		result.SourceIp = sourceIps
	}

	// Extract Managed Device
	managedDeviceValue, ok := conditions.Attributes()["managed_device"]
	if ok && !managedDeviceValue.IsNull() && !managedDeviceValue.IsUnknown() {
		result.ManagedDevice = extractDeviceModel(ctx, managedDeviceValue.(types.Object))
	}

	// Extract Unmanaged Device
	unmanagedDeviceValue, ok := conditions.Attributes()["unmanaged_device"]
	if ok && !unmanagedDeviceValue.IsNull() && !unmanagedDeviceValue.IsUnknown() {
		result.UnmanagedDevice = extractDeviceModel(ctx, unmanagedDeviceValue.(types.Object))
	}

	return result, nil
}

func extractDeviceModel(ctx context.Context, device types.Object) dto.Device {
	result := dto.Device{}

	opswatValue, ok := device.Attributes()["opswat"]
	if ok && !opswatValue.IsNull() && !opswatValue.IsUnknown() {
		result.OpswatMetaAccess = opswatValue.(types.Bool).ValueBool()
	}

	symantecCloudSocValue, ok := device.Attributes()["symantec_cloudsoc"]
	if ok && !symantecCloudSocValue.IsNull() && !symantecCloudSocValue.IsUnknown() {
		result.SymantecCloudSoc = symantecCloudSocValue.(types.Bool).ValueBool()
	}

	symantecWebSecurityServiceValue, ok := device.Attributes()["symantec_web_security_service"]
	if ok && !symantecWebSecurityServiceValue.IsNull() && !symantecWebSecurityServiceValue.IsUnknown() {
		result.SymantecWebSecurityService = symantecWebSecurityServiceValue.(types.Bool).ValueBool()
	}

	return result
}
