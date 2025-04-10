package dto

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestFromPolicyRulesContainers(t *testing.T) {
	tests := []struct {
		name          string
		policyRules   []sdk.PolicyRule
		expectedRules []ActivityRule
	}{
		{
			name:          "Empty policy rules",
			policyRules:   []sdk.PolicyRule{},
			expectedRules: []ActivityRule{},
		},
		{
			name:          "Nil policy rules",
			policyRules:   nil,
			expectedRules: nil,
		},
		{
			name: "Single rule with file downloaded",
			policyRules: []sdk.PolicyRule{
				{
					ActionId: BlockAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: FileDownloadedCondition,
							Arguments:             map[string][]string{},
						},
					},
				},
			},
			expectedRules: []ActivityRule{
				{
					Action: BlockAction,
					Conditions: &RuleConditions{
						FileDownloaded: true,
						Arguments:      &RuleConditionArguments{},
					},
				},
			},
		},
		{
			name: "Single rule with file uploaded",
			policyRules: []sdk.PolicyRule{
				{
					ActionId: BlockUserAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: FileUploadedCondition,
							Arguments:             map[string][]string{},
						},
					},
				},
			},
			expectedRules: []ActivityRule{
				{
					Action: BlockUserAction,
					Conditions: &RuleConditions{
						FileUploaded: true,
						Arguments:    &RuleConditionArguments{},
					},
				},
			},
		},
		{
			name: "Single rule with URI accessed",
			policyRules: []sdk.PolicyRule{
				{
					ActionId: DisconnectUserAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: URICondition,
							Arguments: map[string][]string{
								URIListRuleConditionArgument: {"uri1", "uri2"},
							},
						},
					},
				},
			},
			expectedRules: []ActivityRule{
				{
					Action: DisconnectUserAction,
					Conditions: &RuleConditions{
						UriAccessed: true,
						Arguments: &RuleConditionArguments{
							UriList: []string{"uri1", "uri2"},
						},
					},
				},
			},
		},
		{
			name: "Single rule with HTTP command",
			policyRules: []sdk.PolicyRule{
				{
					ActionId: BlockAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: HTTPCommandCondition,
							Arguments: map[string][]string{
								HTTPCommandRuleConditionArgument: {"GET", "POST"},
							},
						},
					},
				},
			},
			expectedRules: []ActivityRule{
				{
					Action: BlockAction,
					Conditions: &RuleConditions{
						HttpCommand: true,
						Arguments: &RuleConditionArguments{
							Commands: []string{"GET", "POST"},
						},
					},
				},
			},
		},
		{
			name: "Multiple rules with different conditions",
			policyRules: []sdk.PolicyRule{
				{
					ActionId: BlockAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: FileDownloadedCondition,
							Arguments:             map[string][]string{},
						},
					},
				},
				{
					ActionId: BlockUserAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: URICondition,
							Arguments: map[string][]string{
								URIListRuleConditionArgument: {"uri1"},
							},
						},
					},
				},
			},
			expectedRules: []ActivityRule{
				{
					Action: BlockAction,
					Conditions: &RuleConditions{
						FileDownloaded: true,
						Arguments:      &RuleConditionArguments{},
					},
				},
				{
					Action: BlockUserAction,
					Conditions: &RuleConditions{
						UriAccessed: true,
						Arguments: &RuleConditionArguments{
							UriList: []string{"uri1"},
						},
					},
				},
			},
		},
		{
			name: "Web Isolation and DLP actions conditions",
			policyRules: []sdk.PolicyRule{
				{
					ActionId: DLPCloudDetectionAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: FileDownloadedCondition,
							Arguments:             map[string][]string{},
						},
					},
					DlpFilterId: "6fd0a892-8b70-471a-9dd7-bf374b07451f",
				},
				{
					ActionId: IsolationProfile,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: URICondition,
							Arguments: map[string][]string{
								URIListRuleConditionArgument: {"uri1"},
							},
						},
					},
					IsolationProfileId: "571136d7-7bb7-45bc-b039-6e9eea0cc430",
				},
			},
			expectedRules: []ActivityRule{
				{
					Action: DLPCloudDetectionAction,
					Conditions: &RuleConditions{
						FileDownloaded: true,
						Arguments:      &RuleConditionArguments{},
					},
					DLPFilterID: "6fd0a892-8b70-471a-9dd7-bf374b07451f",
				},
				{
					Action: IsolationProfile,
					Conditions: &RuleConditions{
						UriAccessed: true,
						Arguments: &RuleConditionArguments{
							UriList: []string{"uri1"},
						},
					},
					IsolationProfileID: "571136d7-7bb7-45bc-b039-6e9eea0cc430",
				},
			},
		},
		{
			name: "Rule with no conditions",
			policyRules: []sdk.PolicyRule{
				{
					ActionId: BlockAction,
				},
			},
			expectedRules: []ActivityRule{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRules := FromPolicyRulesContainers(tt.policyRules)
			if !assert.Equal(t, actualRules, tt.expectedRules) {
				t.Errorf("FromPolicyRulesContainers() returned unexpected result. \nGot: %+v\nExpected: %+v", actualRules, tt.expectedRules)
			}
		})
	}
}

func TestToPolicyRulesContainers(t *testing.T) {
	tests := []struct {
		name          string
		activityRules []ActivityRule
		expectedRules []sdk.PolicyRule
	}{
		{
			name:          "Empty activity rules",
			activityRules: []ActivityRule{},
			expectedRules: []sdk.PolicyRule{},
		},
		{
			name:          "Nil activity rules",
			activityRules: nil,
			expectedRules: nil,
		},
		{
			name: "Single rule with file downloaded",
			activityRules: []ActivityRule{
				{
					Action: BlockAction,
					Conditions: &RuleConditions{
						FileDownloaded: true,
					},
				},
			},
			expectedRules: []sdk.PolicyRule{
				{
					ActionId: BlockAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: FileDownloadedCondition,
							Arguments:             map[string][]string{},
						},
					},
				},
			},
		},
		{
			name: "Single rule with file uploaded",
			activityRules: []ActivityRule{
				{
					Action: BlockUserAction,
					Conditions: &RuleConditions{
						FileUploaded: true,
					},
				},
			},
			expectedRules: []sdk.PolicyRule{
				{
					ActionId: BlockUserAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: FileUploadedCondition,
							Arguments:             map[string][]string{},
						},
					},
				},
			},
		},
		{
			name: "Single rule with URI accessed",
			activityRules: []ActivityRule{
				{
					Action: DisconnectUserAction,
					Conditions: &RuleConditions{
						UriAccessed: true,
						Arguments: &RuleConditionArguments{
							UriList: []string{"uri1", "uri2"},
						},
					},
				},
			},
			expectedRules: []sdk.PolicyRule{
				{
					ActionId: DisconnectUserAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: URICondition,
							Arguments: map[string][]string{
								URIListRuleConditionArgument: {"uri1", "uri2"},
							},
						},
					},
				},
			},
		},
		{
			name: "Single rule with HTTP command",
			activityRules: []ActivityRule{
				{
					Action: BlockAction,
					Conditions: &RuleConditions{
						HttpCommand: true,
						Arguments: &RuleConditionArguments{
							Commands: []string{"GET", "POST"},
						},
					},
				},
			},
			expectedRules: []sdk.PolicyRule{
				{
					ActionId: BlockAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: HTTPCommandCondition,
							Arguments: map[string][]string{
								HTTPCommandRuleConditionArgument: {"GET", "POST"},
							},
						},
					},
				},
			},
		},
		{
			name: "Multiple rules with different conditions",
			activityRules: []ActivityRule{
				{
					Action: BlockAction,
					Conditions: &RuleConditions{
						FileDownloaded: true,
					},
				},
				{
					Action: BlockUserAction,
					Conditions: &RuleConditions{
						UriAccessed: true,
						Arguments: &RuleConditionArguments{
							UriList: []string{"uri1"},
						},
					},
				},
			},
			expectedRules: []sdk.PolicyRule{
				{
					ActionId: BlockAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: FileDownloadedCondition,
							Arguments:             map[string][]string{},
						},
					},
				},
				{
					ActionId: BlockUserAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: URICondition,
							Arguments: map[string][]string{
								URIListRuleConditionArgument: {"uri1"},
							},
						},
					},
				},
			},
		},
		{
			name: "Web Isolation and DLP actions conditions",
			activityRules: []ActivityRule{
				{
					Action: DLPCloudDetectionAction,
					Conditions: &RuleConditions{
						FileDownloaded: true,
						Arguments:      &RuleConditionArguments{},
					},
					DLPFilterID: "6fd0a892-8b70-471a-9dd7-bf374b07451f",
				},
				{
					Action: IsolationProfile,
					Conditions: &RuleConditions{
						UriAccessed: true,
						Arguments: &RuleConditionArguments{
							UriList: []string{"uri1"},
						},
					},
					IsolationProfileID: "571136d7-7bb7-45bc-b039-6e9eea0cc430",
				},
			},
			expectedRules: []sdk.PolicyRule{
				{
					ActionId: DLPCloudDetectionAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: FileDownloadedCondition,
							Arguments:             map[string][]string{},
						},
					},
					DlpFilterId: "6fd0a892-8b70-471a-9dd7-bf374b07451f",
				},
				{
					ActionId: IsolationProfile,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: URICondition,
							Arguments: map[string][]string{
								URIListRuleConditionArgument: {"uri1"},
							},
						},
					},
					IsolationProfileId: "571136d7-7bb7-45bc-b039-6e9eea0cc430",
				},
			},
		},
		{
			name: "Rule with no conditions",
			activityRules: []ActivityRule{
				{
					Action: BlockAction,
				},
			},
			expectedRules: []sdk.PolicyRule{},
		},
		{
			name: "Rule with empty action",
			activityRules: []ActivityRule{
				{
					Action: "",
					Conditions: &RuleConditions{
						FileDownloaded: true,
					},
				},
			},
			expectedRules: []sdk.PolicyRule{},
		},
		{
			name: "Rule with nil conditions",
			activityRules: []ActivityRule{
				{
					Action:     BlockAction,
					Conditions: nil,
				},
			},
			expectedRules: []sdk.PolicyRule{},
		},
		{
			name: "Multiple conditions",
			activityRules: []ActivityRule{
				{
					Action: BlockAction,
					Conditions: &RuleConditions{
						FileDownloaded: true,
						FileUploaded:   true,
						UriAccessed:    true,
						HttpCommand:    true,
						Arguments: &RuleConditionArguments{
							UriList:  []string{"uri1", "uri2"},
							Commands: []string{"GET", "POST"},
						},
					},
				},
			},
			expectedRules: []sdk.PolicyRule{
				{
					ActionId: BlockAction,
					Conditions: []sdk.PolicyCondition{
						{
							ConditionDefinitionId: FileDownloadedCondition,
							Arguments:             map[string][]string{},
						},
						{
							ConditionDefinitionId: FileUploadedCondition,
							Arguments:             map[string][]string{},
						},
						{
							ConditionDefinitionId: URICondition,
							Arguments: map[string][]string{
								URIListRuleConditionArgument: {"uri1", "uri2"},
							},
						},
						{
							ConditionDefinitionId: HTTPCommandCondition,
							Arguments: map[string][]string{
								HTTPCommandRuleConditionArgument: {"GET", "POST"},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRules := ToPolicyRulesContainers(tt.activityRules)
			if !reflect.DeepEqual(actualRules, tt.expectedRules) {
				t.Errorf("ToPolicyRulesContainers() returned unexpected result. \nGot: %+v\nExpected: %+v", actualRules, tt.expectedRules)
			}
		})
	}
}
