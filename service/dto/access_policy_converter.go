package dto

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"errors"
	"fmt"
)

func FromTargetProtocol(targetProtocol sdk.PolicyTargetProtocol) string {
	switch targetProtocol {
	case sdk.HTTP_PolicyTargetProtocol:
		return "HTTP"
	case sdk.SSH_PolicyTargetProtocol:
		return "SSH"
	case sdk.RDP_PolicyTargetProtocol:
		return "RDP"
	case sdk.TCP_PolicyTargetProtocol:
		return "TCP"
	default:
		return ""
	}
}

func ToTargetProtocol(targetProtocol string) *sdk.PolicyTargetProtocol {
	var policyTargetProtocol sdk.PolicyTargetProtocol

	switch targetProtocol {
	case "HTTP":
		policyTargetProtocol = sdk.HTTP_PolicyTargetProtocol
	case "SSH":
		policyTargetProtocol = sdk.SSH_PolicyTargetProtocol
	case "RDP":
		policyTargetProtocol = sdk.RDP_PolicyTargetProtocol
	case "TCP":
		policyTargetProtocol = sdk.TCP_PolicyTargetProtocol
	}

	return &policyTargetProtocol
}

func ToApplicationType(targetProtocol string) *sdk.ApplicationType {
	var applicationType sdk.ApplicationType

	switch targetProtocol {
	case "HTTP":
		applicationType = sdk.HTTP_ApplicationType
	case "SSH":
		applicationType = sdk.SSH_ApplicationType
	case "RDP":
		applicationType = sdk.RDP_ApplicationType
	case "TCP":
		applicationType = sdk.TCP_ApplicationType
	}

	return &applicationType
}

func FromModelType(modelType sdk.EntityType) string {

	switch modelType {
	case sdk.USER_EntityType:
		return "User"
	case sdk.GROUP_EntityType:
		return "Group"
	case sdk.API_CLIENT_EntityType:
		return "ApiClient"
	}

	return ""
}

func ToModelType(entityType string) *sdk.EntityType {
	var modelType sdk.EntityType

	switch entityType {
	case "User":
		modelType = sdk.USER_EntityType
	case "Group":
		modelType = sdk.GROUP_EntityType
	case "ApiClient":
		modelType = sdk.API_CLIENT_EntityType
	}

	return &modelType
}

func ConvertIdentityProviderTypeToString(idpType interface{}) string {

	if idpType == "" || idpType == nil {
		return ""
	}
	return fmt.Sprintf("%s", idpType)
}

func ConvertIdentityProviderTypeToEnum(idpType string) (sdk.IdentityProviderType, error) {
	switch idpType {
	case "local", "keycloak":
		return sdk.LOCAL_IdentityProviderType, nil
	case "ad", "azuread": //PLT-117 - ad and azuread are synonyms - referring to Azure AD.
		return sdk.AD_IdentityProviderType, nil
	case "okta":
		return sdk.OKTA_IdentityProviderType, nil
	case "adfs":
		return sdk.ADFS_IdentityProviderType, nil
	case "gapps":
		return sdk.GAPPS_IdentityProviderType, nil
	case "onelogin":
		return sdk.ONELOGIN_IdentityProviderType, nil
	}
	return "", errors.New("Failed to locate matching provider type")
}

func ConvertToDto(accessPolicy *AccessPolicy) sdk.PolicyAccess {
	accessPolicyType := sdk.ACCESS_PolicyType

	var rdpSettingsDto *sdk.PolicyRdpSettings
	var sshSettingsDto *sdk.PolicySshSettings
	var tcpSettingsDto *sdk.PolicyTcpSettings
	var directoryEntities []sdk.DirectoryEntity
	var applications []sdk.ApplicationBase
	var validatorsDto map[string]bool
	var conditionsDto []sdk.PolicyCondition

	directoryEntities = EntityDTOToEntityModel(accessPolicy.DirectoryEntities)

	for _, applicationId := range accessPolicy.Applications {
		applications = append(applications, sdk.ApplicationBase{
			Id:    applicationId,
			Type_: ToApplicationType(accessPolicy.TargetProtocol),
		})
	}

	conditionsDto = ToFilterConditions(accessPolicy.Conditions)

	if accessPolicy.RdpSettings != nil {
		rdpSettingsDto = &sdk.PolicyRdpSettings{
			LongTermPassword: accessPolicy.RdpSettings.LongTermPassword,
		}
	}

	if accessPolicy.SshSettings != nil {
		sshSettingsDto = &sdk.PolicySshSettings{
			Accounts:             accessPolicy.SshSettings.Accounts,
			AutoMapping:          accessPolicy.SshSettings.AutoMapping,
			FullUpnAutoMapping:   accessPolicy.SshSettings.FullUPNAutoMapping,
			AgentForward:         accessPolicy.SshSettings.AgentForward,
			AcceptTemporaryToken: accessPolicy.SshSettings.AcceptTemporaryToken,
			AcceptCertificate:    accessPolicy.SshSettings.AcceptCertificate,
		}
	}

	if accessPolicy.TcpSettings != nil {
		tcpSettingsDto = &sdk.PolicyTcpSettings{
			AcceptTemporaryToken: accessPolicy.TcpSettings.AcceptTemporaryToken,
			AcceptCertificate:    accessPolicy.TcpSettings.AcceptCertificate,
		}
	}

	if accessPolicy.Validators != nil {
		validatorsDto = map[string]bool{}

		if accessPolicy.Validators.WebVerification {
			validatorsDto[ValidatorWebVerification] = accessPolicy.Validators.WebVerification
		}
		if accessPolicy.Validators.MFA {
			validatorsDto[MFA] = accessPolicy.Validators.MFA
		}
	}

	accessPolicyDto := sdk.PolicyAccess{
		Type_:             &accessPolicyType,
		TargetProtocol:    ToTargetProtocol(accessPolicy.TargetProtocol),
		Id:                accessPolicy.Id,
		CollectionId:      accessPolicy.CollectionID,
		Enabled:           accessPolicy.Enabled,
		CreatedAt:         accessPolicy.CreatedAt,
		Name:              accessPolicy.Name,
		DirectoryEntities: directoryEntities,
		Applications:      applications,
		FilterConditions:  conditionsDto,
		Validators:        validatorsDto,
		RdpSettings:       rdpSettingsDto,
		SshSettings:       sshSettingsDto,
		TcpSettings:       tcpSettingsDto,
	}

	return accessPolicyDto
}

func ToFilterConditions(conditions *Conditions) []sdk.PolicyCondition {
	var conditionsDto []sdk.PolicyCondition
	if conditions != nil {
		if conditions.SourceIp != nil {
			conditionsDto = append(conditionsDto, sdk.PolicyCondition{
				ConditionDefinitionId: IpCondition,
				Arguments: map[string][]string{
					IpUuid:           conditions.SourceIp,
					SharedIpListUuid: conditions.SharedIpList,
				},
			})
		}

		if conditions.Location != nil {
			conditionsDto = append(conditionsDto, sdk.PolicyCondition{
				ConditionDefinitionId: LocationRestrictionCondition,
				Arguments: map[string][]string{
					CountriesUuid: conditions.Location,
				},
			})
		}

		if hasDeviceArgument(conditions.ManagedDevice) {
			managedArgumentsMap := getDeviceArguments(conditions.ManagedDevice)
			conditionsDto = append(conditionsDto, sdk.PolicyCondition{
				ConditionDefinitionId: ManagedDeviceCondition,
				Arguments:             managedArgumentsMap,
			})
		}

		if hasDeviceArgument(conditions.UnmanagedDevice) {
			unmanagedArgumentsMap := getDeviceArguments(conditions.UnmanagedDevice)
			conditionsDto = append(conditionsDto, sdk.PolicyCondition{
				ConditionDefinitionId: UnmanagedDeviceCondition,
				Arguments:             unmanagedArgumentsMap,
			})
		}
	}
	return conditionsDto
}

func getDeviceArguments(deviceArg Device) map[string][]string {
	var deviceArguments []string

	deviceArgumentsMap := map[string][]string{}

	if deviceArg.SymantecWebSecurityService {
		deviceArguments = append(deviceArguments, ManagedDeviceWssConditionArgument)
	}

	if deviceArg.SymantecCloudSoc {
		deviceArguments = append(deviceArguments, ManagedDeviceCloudSocConditionArgument)
	}

	if deviceArg.OpswatMetaAccess {
		deviceArguments = append(deviceArguments, ManagedDeviceOpswatConditionArgument)
	}

	if deviceArguments != nil {
		deviceArgumentsMap[Authentication] = deviceArguments
	}

	return deviceArgumentsMap
}

func hasDeviceArgument(deviceArg Device) bool {
	if deviceArg.OpswatMetaAccess || deviceArg.SymantecCloudSoc || deviceArg.SymantecWebSecurityService {
		return true
	}
	return false
}

func ConvertFromDto(accessPolicyDto sdk.PolicyAccess) *AccessPolicy {
	var applications []string
	var directoryEntity []DirectoryEntity
	var rdpSetting *PolicyRdpSettings
	var sshSetting *PolicySshSettings
	var tcpSetting *PolicyTcpSettings
	var validators *Validators
	var conditions *Conditions

	for _, applicationsDto := range accessPolicyDto.Applications {
		applications = append(applications, applicationsDto.Id)
	}

	directoryEntity = EntityModelEntityDTO(accessPolicyDto.DirectoryEntities)

	conditions = FromFilterConditions(accessPolicyDto.FilterConditions)

	if accessPolicyDto.RdpSettings != nil {
		rdpSetting = &PolicyRdpSettings{
			LongTermPassword: accessPolicyDto.RdpSettings.LongTermPassword,
		}
	}

	if accessPolicyDto.SshSettings != nil {
		sshSetting = &PolicySshSettings{
			Accounts:             accessPolicyDto.SshSettings.Accounts,
			AutoMapping:          accessPolicyDto.SshSettings.AutoMapping,
			FullUPNAutoMapping:   accessPolicyDto.SshSettings.FullUpnAutoMapping,
			AgentForward:         accessPolicyDto.SshSettings.AgentForward,
			AcceptTemporaryToken: accessPolicyDto.SshSettings.AcceptTemporaryToken,
			AcceptCertificate:    accessPolicyDto.SshSettings.AcceptCertificate,
		}
	}

	if accessPolicyDto.TcpSettings != nil {
		tcpSetting = &PolicyTcpSettings{
			AcceptTemporaryToken: accessPolicyDto.TcpSettings.AcceptTemporaryToken,
			AcceptCertificate:    accessPolicyDto.TcpSettings.AcceptCertificate,
		}
	}

	if accessPolicyDto.Validators != nil && len(accessPolicyDto.Validators) > 0 {
		validators = &Validators{
			WebVerification: accessPolicyDto.Validators[ValidatorWebVerification],
			MFA:             accessPolicyDto.Validators[MFA],
		}
	}

	return &AccessPolicy{
		Policy: Policy{
			TargetProtocol:    FromTargetProtocol(*accessPolicyDto.TargetProtocol),
			CollectionID:      accessPolicyDto.CollectionId,
			Id:                accessPolicyDto.Id,
			Enabled:           accessPolicyDto.Enabled,
			CreatedAt:         accessPolicyDto.CreatedAt,
			Name:              accessPolicyDto.Name,
			DirectoryEntities: directoryEntity,
			Applications:      applications,
			Conditions:        conditions,
		},
		Validators:  validators,
		RdpSettings: rdpSetting,
		SshSettings: sshSetting,
		TcpSettings: tcpSetting,
	}
}

func FromFilterConditions(filterConditions []sdk.PolicyCondition) *Conditions {
	var conditions *Conditions
	if filterConditions != nil && len(filterConditions) > 0 {
		conditions = &Conditions{}

		for _, filterCondition := range filterConditions {
			if filterCondition.ConditionDefinitionId == IpCondition {
				for _, ipCondition := range filterCondition.Arguments[IpUuid] {
					conditions.SourceIp = append(conditions.SourceIp, ipCondition)
				}

				if _, ok := filterCondition.Arguments[SharedIpListUuid]; ok {
					for _, sharedIpListCondition := range filterCondition.Arguments[SharedIpListUuid] {
						conditions.SharedIpList = append(conditions.SharedIpList, sharedIpListCondition)
					}
				}
			}

			if filterCondition.ConditionDefinitionId == LocationRestrictionCondition {
				for _, locationCondition := range filterCondition.Arguments[CountriesUuid] {
					conditions.Location = append(conditions.Location, locationCondition)
				}
			}

			if filterCondition.ConditionDefinitionId == ManagedDeviceCondition || filterCondition.ConditionDefinitionId == UnmanagedDeviceCondition {
				for _, deviceCondition := range filterCondition.Arguments[Authentication] {

					if ManagedDeviceOpswatConditionArgument == deviceCondition {
						conditions.ManagedDevice.OpswatMetaAccess = true
					}

					if ManagedDeviceCloudSocConditionArgument == deviceCondition {
						conditions.ManagedDevice.SymantecCloudSoc = true
					}

					if ManagedDeviceWssConditionArgument == deviceCondition {
						conditions.ManagedDevice.SymantecWebSecurityService = true
					}
				}
			}

		}
	}
	return conditions
}
