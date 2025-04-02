package framework_provider

import (
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

func NewWebActivityPolicyResource() resource.Resource {
	return &WebActivityPolicyResource{}
}

type WebActivityPolicyResource struct {
	BasePolicyResource
}

type WebActivityPolicyResourceModel struct {
	BasePolicyResourceModel
	Rules types.List `tfsdk:"rules"`
}

func (w *WebActivityPolicyResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_web_activity_policy"
}

func (w *WebActivityPolicyResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	policyAttributes := CreatePolicyBaseSchemaAttributes()

	policyAttributes["rules"] = schema.ListNestedAttribute{
		Required: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"action": schema.StringAttribute{
					Required:    true,
					Description: "the action to apply for this rule condition.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							dto.BlockAction,
							dto.BlockUserAction,
							dto.DisconnectUserAction,
						),
					},
				},
				"conditions": schema.SingleNestedAttribute{
					Optional:    true,
					Description: "the rule conditions arguments if required per enabled condition",
					Attributes: map[string]schema.Attribute{
						"file_downloaded": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
							Description: "File Downloaded rule condition enabled",
						},
						"file_uploaded": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
							Description: "File Uploaded rule condition enabled",
						},
						"uri_accessed": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
							Description: "URI Accessed rule condition enabled",
						},
						"http_command": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
							Description: "HTTP Command rule enabled",
						},
						"arguments": schema.SingleNestedAttribute{
							Optional:    true,
							Computed:    true,
							Description: "the rule conditions arguments if required per enabled condition",
							Attributes: map[string]schema.Attribute{
								"uri_list": schema.ListAttribute{
									Description: "the list of URI to apply URI Accessed rule condition.",
									Optional:    true,
									Computed:    true,
									ElementType: types.StringType,
								},
								"commands": schema.ListAttribute{
									Description: "the HTTP commands to apply HTTP Command rule condition.",
									Optional:    true,
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
		},
	}

	response.Schema = schema.Schema{
		MarkdownDescription: "Web activity policy resource",
		Attributes:          policyAttributes,
	}
}

func (r *WebActivityPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (w *WebActivityPolicyResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data WebActivityPolicyResourceModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	policy, diags := extractActivityPolicyDTO(ctx, &data)

	if diags != nil && diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	for i := range policy.DirectoryEntities {
		resolvedIdentityProviderType, err := w.client.IdentityProviders.GetIdentityProviderTypeById(policy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			response.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to lookup identity provider type for identity provider id %s, got error: %s", policy.DirectoryEntities[i].IdentityProviderId, err))
			return
		}
		policy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)

		// Get Display Name for User/Group by ID
		var resolvedDisplayName string
		switch strings.ToLower(policy.DirectoryEntities[i].EntityType) {
		case "user":
			resolvedDisplayName, err = w.client.IdentityProviders.GetUserDisplayNameTypeById(policy.DirectoryEntities[i].IdentityProviderId, policy.DirectoryEntities[i].IdentifierInProvider)
		case "group":
			resolvedDisplayName, err = w.client.IdentityProviders.GetGroupDisplayNameTypeById(policy.DirectoryEntities[i].IdentityProviderId, policy.DirectoryEntities[i].IdentifierInProvider)
		default:
			response.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to lookup displayName - unknown entity type \"%s\"", policy.DirectoryEntities[i].EntityType))
			return
		}

		if err != nil {
			response.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to lookup displayName for entity type %s with identifier id %s on Identity Provider ID %s", policy.DirectoryEntities[i].EntityType, policy.DirectoryEntities[i].IdentifierInProvider, policy.DirectoryEntities[i].IdentityProviderId))
			return
		}
		policy.DirectoryEntities[i].DisplayName = resolvedDisplayName
	}

	createdPolicy, err := w.client.ActivityPolicies.CreateActivityPolicy(policy)
	if err != nil {
		response.Diagnostics.AddError("Service Error", fmt.Sprintf("Unable to create web activity policy, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a web activity policy resource ID="+createdPolicy.Id)

	model, diags := w.readWebActivityPolicy(ctx, createdPolicy.Id)
	if diags != nil && diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	if model == nil {
		response.Diagnostics.AddError("Service Error", "Failed to get web activity policy")
		return
	}

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, model)...)
}

func (w *WebActivityPolicyResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data WebActivityPolicyResourceModel

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(response.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	model, diags := w.readWebActivityPolicy(ctx, data.ID.ValueString())
	if diags != nil && diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	// if remove from state if it does not exist
	if model == nil {
		response.Diagnostics.AddWarning("Service Warning", "Web activity policy does not exist, removing from state")
		response.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, model)...)
}

func (w *WebActivityPolicyResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var currentData WebActivityPolicyResourceModel
	var data WebActivityPolicyResourceModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	// Read Terraform state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &currentData)...)

	if response.Diagnostics.HasError() {
		return
	}

	// Use the ID from the state
	data.ID = currentData.ID

	policy, diags := extractActivityPolicyDTO(ctx, &data)

	if diags != nil && diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	for i := range policy.DirectoryEntities {
		resolvedIdentityProviderType, err := w.client.IdentityProviders.GetIdentityProviderTypeById(policy.DirectoryEntities[i].IdentityProviderId)
		if err != nil {
			response.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to lookup identity provider type for identity provider id %s, got error: %s", policy.DirectoryEntities[i].IdentityProviderId, err))
			return
		}
		policy.DirectoryEntities[i].IdentityProviderType = dto.ConvertIdentityProviderTypeToString(resolvedIdentityProviderType)
	}

	updatedPolicy, err := w.client.ActivityPolicies.UpdateActivityPolicy(policy)
	if err != nil {
		response.Diagnostics.AddError("Service Error", fmt.Sprintf("Unable to update web activity policy, got error: %s", err))
		return
	}

	model, diags := w.readWebActivityPolicy(ctx, updatedPolicy.Id)
	if diags != nil && diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	if model == nil {
		response.Diagnostics.AddError("Service Error", "Failed to get web activity policy")
		return
	}

	tflog.Trace(ctx, "created a web activity policy resource ID="+updatedPolicy.Id)

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, model)...)
}

func (w *WebActivityPolicyResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data WebActivityPolicyResourceModel

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	err := w.client.ActivityPolicies.DeleteActivityPolicy(data.ID.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Service Error", fmt.Sprintf("Unable to delete web activity policy, got error: %s", err))
		return
	}
}

func (w *WebActivityPolicyResource) readWebActivityPolicy(ctx context.Context, policyID string) (*WebActivityPolicyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	activityPolicy, err := w.client.ActivityPolicies.GetActivityPolicy(policyID)
	if err != nil {
		diags.AddError("Service Error", fmt.Sprintf("Unable to read web activity policy, got error: %s", err))
		return nil, diags
	}

	if activityPolicy == nil {
		return nil, nil
	}

	model, diags := extractActivityPolicyModel(ctx, activityPolicy)
	if diags != nil && diags.HasError() {
		return nil, diags
	}
	return model, nil
}

func extractActivityPolicyDTO(ctx context.Context, webActivityModel *WebActivityPolicyResourceModel) (*dto.ActivityPolicy, diag.Diagnostics) {
	policy, diags := convertBaseModelToPolicy(ctx, &webActivityModel.BasePolicyResourceModel)
	if diags != nil && diags.HasError() {
		return nil, diags
	}
	policy.TargetProtocol = "HTTP"

	activityRules, diags := extractActivityRules(ctx, webActivityModel)
	if diags != nil && diags.HasError() {
		return nil, diags
	}

	webActivityPolicy := &dto.ActivityPolicy{
		Policy:        *policy,
		ActivityRules: activityRules,
	}
	return webActivityPolicy, nil
}

func extractActivityPolicyModel(ctx context.Context, activityPolicy *dto.ActivityPolicy) (*WebActivityPolicyResourceModel, diag.Diagnostics) {
	policyModel, diags := convertPolicyToBaseModel(ctx, &activityPolicy.Policy)
	if diags != nil && diags.HasError() {
		return nil, diags
	}

	activityPolicyModelRules, diags := flattenActivityRules(ctx, activityPolicy.ActivityRules)
	if diags != nil && diags.HasError() {
		return nil, diags
	}

	webActivityPolicyModel := &WebActivityPolicyResourceModel{
		BasePolicyResourceModel: *policyModel,
		Rules:                   activityPolicyModelRules,
	}
	return webActivityPolicyModel, nil
}

func extractActivityRules(ctx context.Context, webActivityModel *WebActivityPolicyResourceModel) ([]dto.ActivityRule, diag.Diagnostics) {
	var activityRules []dto.ActivityRule
	var diags diag.Diagnostics

	if webActivityModel.Rules.IsNull() || webActivityModel.Rules.IsUnknown() {
		return activityRules, diags
	}

	var rules []types.Object
	diags = webActivityModel.Rules.ElementsAs(ctx, &rules, false)
	if diags.HasError() {
		return activityRules, diags
	}

	for _, rule := range rules {
		// Extract Action
		actionValue, ok := rule.Attributes()["action"]
		if !ok || actionValue.IsNull() || actionValue.IsUnknown() {
			diags.AddError("Missing Rule Action", "The action attribute is required for each rule.")
			return activityRules, diags
		}
		action := actionValue.(types.String).ValueString()

		conditions, diagnostics := extractRuleConditions(ctx, rule)
		if diagnostics != nil && diagnostics.HasError() {
			return nil, diagnostics
		}

		activityRule := dto.ActivityRule{
			Action:     action,
			Conditions: conditions,
		}
		activityRules = append(activityRules, activityRule)
	}

	return activityRules, diags
}

func extractRuleConditions(ctx context.Context, rule types.Object) (*dto.RuleConditions, diag.Diagnostics) {
	conditions := dto.RuleConditions{}

	conditionsValue, ok := rule.Attributes()["conditions"]
	if !ok || conditionsValue.IsNull() || conditionsValue.IsUnknown() {
		return &conditions, nil
	}
	conditionsObject := conditionsValue.(types.Object)

	// File Downloaded
	fileDownloadedValue, ok := conditionsObject.Attributes()["file_downloaded"]
	if ok && !fileDownloadedValue.IsNull() && !fileDownloadedValue.IsUnknown() {
		conditions.FileDownloaded = fileDownloadedValue.(types.Bool).ValueBool()
	}

	// File Uploaded
	fileUploadedValue, ok := conditionsObject.Attributes()["file_uploaded"]
	if ok && !fileUploadedValue.IsNull() && !fileUploadedValue.IsUnknown() {
		conditions.FileUploaded = fileUploadedValue.(types.Bool).ValueBool()
	}

	// URI Accessed
	uriAccessedValue, ok := conditionsObject.Attributes()["uri_accessed"]
	if ok && !uriAccessedValue.IsNull() && !uriAccessedValue.IsUnknown() {
		conditions.UriAccessed = uriAccessedValue.(types.Bool).ValueBool()
	}

	// HTTP Command
	httpCommandValue, ok := conditionsObject.Attributes()["http_command"]
	if ok && !httpCommandValue.IsNull() && !httpCommandValue.IsUnknown() {
		conditions.HttpCommand = httpCommandValue.(types.Bool).ValueBool()
	}

	// Condition Rule Arguments
	arguments, argumentsDiags := extractRuleArguments(ctx, conditionsObject)
	if argumentsDiags != nil && argumentsDiags.HasError() {
		return nil, argumentsDiags
	}
	conditions.Arguments = arguments
	return &conditions, nil
}

func extractRuleArguments(ctx context.Context, condition types.Object) (*dto.RuleConditionArguments, diag.Diagnostics) {
	var diags diag.Diagnostics
	arguments := &dto.RuleConditionArguments{}

	argumentsValue, ok := condition.Attributes()["arguments"]
	if ok && !argumentsValue.IsNull() && !argumentsValue.IsUnknown() {
		argumentsObject := argumentsValue.(types.Object)

		// Extract URI List
		uriListValue, ok := argumentsObject.Attributes()["uri_list"]
		if ok && !uriListValue.IsNull() && !uriListValue.IsUnknown() {
			diags.Append(uriListValue.(types.List).ElementsAs(ctx, &arguments.UriList, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		// Extract Commands
		commandsValue, ok := argumentsObject.Attributes()["commands"]
		if ok && !commandsValue.IsNull() && !commandsValue.IsUnknown() {
			diags.Append(commandsValue.(types.List).ElementsAs(ctx, &arguments.Commands, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}
	}

	return arguments, diags
}

func flattenActivityRules(ctx context.Context, activityRules []dto.ActivityRule) (types.List, diag.Diagnostics) {
	emptyRules := types.ListNull(types.ObjectType{AttrTypes: activityRuleAttributeTypes()})
	if len(activityRules) == 0 {
		return emptyRules, nil
	}

	ruleObjects := make([]types.Object, len(activityRules))
	for i, rule := range activityRules {
		ruleObject, ruleDiags := flattenActivityRule(ctx, rule)
		if ruleDiags.HasError() {
			return emptyRules, ruleDiags
		}
		ruleObjects[i] = ruleObject
	}

	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: activityRuleAttributeTypes()}, ruleObjects)
}

func flattenActivityRule(ctx context.Context, activityRule dto.ActivityRule) (types.Object, diag.Diagnostics) {
	ruleAttributes := make(map[string]attr.Value)

	ruleAttributes["action"] = types.StringValue(activityRule.Action)

	conditions, conditionsDiags := flattenRuleConditions(ctx, activityRule.Conditions)
	if conditionsDiags.HasError() {
		return types.ObjectNull(activityRuleAttributeTypes()), conditionsDiags
	}
	ruleAttributes["conditions"] = conditions

	return types.ObjectValue(activityRuleAttributeTypes(), ruleAttributes)
}

func flattenRuleConditions(ctx context.Context, conditions *dto.RuleConditions) (types.Object, diag.Diagnostics) {
	conditionsAttributes := make(map[string]attr.Value)

	if conditions == nil {
		return types.ObjectNull(ruleConditionsAttributeTypes()), nil
	}

	conditionsAttributes["file_downloaded"] = types.BoolValue(conditions.FileDownloaded)
	conditionsAttributes["file_uploaded"] = types.BoolValue(conditions.FileUploaded)
	conditionsAttributes["uri_accessed"] = types.BoolValue(conditions.UriAccessed)
	conditionsAttributes["http_command"] = types.BoolValue(conditions.HttpCommand)

	arguments, argumentsDiags := flattenRuleArguments(ctx, conditions.Arguments)
	if argumentsDiags.HasError() {
		return types.ObjectNull(ruleConditionsAttributeTypes()), argumentsDiags
	}
	conditionsAttributes["arguments"] = arguments

	return types.ObjectValue(ruleConditionsAttributeTypes(), conditionsAttributes)
}

func flattenRuleArguments(ctx context.Context, arguments *dto.RuleConditionArguments) (types.Object, diag.Diagnostics) {
	emptyRuleArguments := types.ObjectNull(conditionsArgumentsAttributeTypes())

	if arguments == nil {
		return emptyRuleArguments, nil
	}

	if len(arguments.UriList) == 0 && len(arguments.Commands) == 0 {
		return emptyRuleArguments, nil

	}

	argumentsAttributes := make(map[string]attr.Value)

	argumentsAttributes["uri_list"] = types.ListNull(types.StringType)
	if len(arguments.UriList) > 0 {
		uriList, uriListDiags := types.ListValueFrom(ctx, types.StringType, arguments.UriList)
		if uriListDiags.HasError() {
			return emptyRuleArguments, uriListDiags
		}
		argumentsAttributes["uri_list"] = uriList
	}

	argumentsAttributes["commands"] = types.ListNull(types.StringType)
	if len(arguments.Commands) > 0 {
		commands, commandsDiags := types.ListValueFrom(ctx, types.StringType, arguments.Commands)
		if commandsDiags.HasError() {
			return emptyRuleArguments, commandsDiags
		}
		argumentsAttributes["commands"] = commands
	}

	return types.ObjectValue(conditionsArgumentsAttributeTypes(), argumentsAttributes)
}

func activityRuleAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"action":     types.StringType,
		"conditions": types.ObjectType{AttrTypes: ruleConditionsAttributeTypes()},
	}
}

func ruleConditionsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"file_downloaded": types.BoolType,
		"file_uploaded":   types.BoolType,
		"uri_accessed":    types.BoolType,
		"http_command":    types.BoolType,
		"arguments":       types.ObjectType{AttrTypes: conditionsArgumentsAttributeTypes()},
	}
}

func conditionsArgumentsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"uri_list": types.ListType{ElemType: types.StringType},
		"commands": types.ListType{ElemType: types.StringType},
	}
}
