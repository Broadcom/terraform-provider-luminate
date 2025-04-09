package dto

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
)

func ToActivityPolicyDto(activityPolicy *ActivityPolicy) sdk.PolicyActivity {
	activityPolicyType := sdk.ACTIVITY_PolicyType

	var directoryEntities []sdk.DirectoryEntity
	var applications []sdk.ApplicationBase
	var conditionsDto []sdk.PolicyCondition
	var policyRules []sdk.PolicyRule

	directoryEntities = EntityDTOToEntityModel(activityPolicy.DirectoryEntities)

	for _, applicationId := range activityPolicy.Applications {
		applications = append(applications, sdk.ApplicationBase{
			Id:    applicationId,
			Type_: ToApplicationType(activityPolicy.TargetProtocol),
		})
	}

	conditionsDto = ToFilterConditions(activityPolicy.Conditions)

	policyRules = ToPolicyRulesContainers(activityPolicy.ActivityRules)

	activityPolicyDto := sdk.PolicyActivity{
		Type_:             &activityPolicyType,
		TargetProtocol:    ToTargetProtocol(activityPolicy.TargetProtocol),
		Id:                activityPolicy.Id,
		CollectionId:      activityPolicy.CollectionID,
		Enabled:           activityPolicy.Enabled,
		CreatedAt:         activityPolicy.CreatedAt,
		Name:              activityPolicy.Name,
		DirectoryEntities: directoryEntities,
		Applications:      applications,
		FilterConditions:  conditionsDto,
		Containers:        policyRules,
		IsIsolation:       activityPolicy.EnableIsolation,
	}

	return activityPolicyDto
}

func FromActivityPolicyDto(activityPolicyDto sdk.PolicyActivity) *ActivityPolicy {
	var applications []string
	var directoryEntity []DirectoryEntity
	var conditions *Conditions
	var activityRules []ActivityRule

	for _, applicationsDto := range activityPolicyDto.Applications {
		applications = append(applications, applicationsDto.Id)
	}

	directoryEntity = EntityModelEntityDTO(activityPolicyDto.DirectoryEntities)

	conditions = FromFilterConditions(activityPolicyDto.FilterConditions)

	activityRules = FromPolicyRulesContainers(activityPolicyDto.Containers)

	activityPolicy := &ActivityPolicy{
		Policy: Policy{
			TargetProtocol:    FromTargetProtocol(*activityPolicyDto.TargetProtocol),
			CollectionID:      activityPolicyDto.CollectionId,
			Id:                activityPolicyDto.Id,
			Enabled:           activityPolicyDto.Enabled,
			CreatedAt:         activityPolicyDto.CreatedAt,
			Name:              activityPolicyDto.Name,
			DirectoryEntities: directoryEntity,
			Applications:      applications,
			Conditions:        conditions,
		},
		ActivityRules:   activityRules,
		EnableIsolation: activityPolicyDto.IsIsolation,
	}
	return activityPolicy
}

func FromPolicyRulesContainers(policyRules []sdk.PolicyRule) []ActivityRule {
	var activityRules []ActivityRule
	if policyRules != nil {
		activityRules = make([]ActivityRule, 0, len(policyRules))
		for _, policyRule := range policyRules {
			if policyRule.Conditions != nil && len(policyRule.Conditions) > 0 {
				ruleConditions := RuleConditions{Arguments: &RuleConditionArguments{}}
				for _, policyRuleCondition := range policyRule.Conditions {
					if policyRuleCondition.ConditionDefinitionId == FileDownloadedCondition {
						ruleConditions.FileDownloaded = true
					}
					if policyRuleCondition.ConditionDefinitionId == FileUploadedCondition {
						ruleConditions.FileUploaded = true
					}
					if policyRuleCondition.ConditionDefinitionId == URICondition &&
						len(policyRuleCondition.Arguments) > 0 {

						uriList, ok := policyRuleCondition.Arguments[URIListRuleConditionArgument]
						if ok {
							ruleConditions.UriAccessed = true
							ruleConditions.Arguments.UriList = uriList
						}
					}
					if policyRuleCondition.ConditionDefinitionId == HTTPCommandCondition &&
						len(policyRuleCondition.Arguments) > 0 {

						commands, ok := policyRuleCondition.Arguments[HTTPCommandRuleConditionArgument]
						if ok {
							ruleConditions.HttpCommand = true
							ruleConditions.Arguments.Commands = commands
						}
					}
				}
				activityRule := ActivityRule{
					Action:             policyRule.ActionId,
					Conditions:         &ruleConditions,
					IsolationProfileID: policyRule.IsolationProfileId,
				}
				activityRules = append(activityRules, activityRule)
			}
		}
	}
	return activityRules
}

func ToPolicyRulesContainers(activityRules []ActivityRule) []sdk.PolicyRule {
	var policyRules []sdk.PolicyRule
	if activityRules != nil {
		policyRules = make([]sdk.PolicyRule, 0, len(activityRules))
		for _, activityRule := range activityRules {
			if activityRule.Conditions == nil || activityRule.Action == "" {
				continue
			}
			policyRule := sdk.PolicyRule{
				ActionId:           activityRule.Action,
				IsolationProfileId: activityRule.IsolationProfileID,
			}
			if activityRule.Conditions.FileDownloaded {
				condition := sdk.PolicyCondition{
					ConditionDefinitionId: FileDownloadedCondition,
					Arguments:             map[string][]string{},
				}
				policyRule.Conditions = append(policyRule.Conditions, condition)
			}
			if activityRule.Conditions.FileUploaded {
				condition := sdk.PolicyCondition{
					ConditionDefinitionId: FileUploadedCondition,
					Arguments:             map[string][]string{},
				}
				policyRule.Conditions = append(policyRule.Conditions, condition)
			}
			if activityRule.Conditions.UriAccessed &&
				activityRule.Conditions.Arguments != nil &&
				activityRule.Conditions.Arguments.UriList != nil {

				condition := sdk.PolicyCondition{
					ConditionDefinitionId: URICondition,
					Arguments:             map[string][]string{URIListRuleConditionArgument: activityRule.Conditions.Arguments.UriList},
				}
				policyRule.Conditions = append(policyRule.Conditions, condition)
			}
			if activityRule.Conditions.HttpCommand &&
				activityRule.Conditions.Arguments != nil &&
				activityRule.Conditions.Arguments.Commands != nil {

				condition := sdk.PolicyCondition{
					ConditionDefinitionId: HTTPCommandCondition,
					Arguments:             map[string][]string{HTTPCommandRuleConditionArgument: activityRule.Conditions.Arguments.Commands},
				}
				policyRule.Conditions = append(policyRule.Conditions, condition)
			}
			policyRules = append(policyRules, policyRule)
		}

	}
	return policyRules
}
