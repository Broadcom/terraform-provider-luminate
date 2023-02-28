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

func FromModelType(modelType string) string {

	switch modelType {
	case "User":
		return "User"
	case "Group":
		return "Group"
	case "ApiClient":
		return "ApiClient"
	}

	return ""
}

func ToModelType(entityType string) *string {
	var modelType string

	switch entityType {
	case "User":
		modelType = "User"
	case "Group":
		modelType = "Group"
	case "ApiClient":
		modelType = "ApiClient"
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

	for _, directoryEntity := range accessPolicy.DirectoryEntities {
		directoryEntities = append(directoryEntities, sdk.DirectoryEntity{
			IdentifierInProvider: directoryEntity.IdentifierInProvider,
			IdentityProviderId:   directoryEntity.IdentityProviderId,
			DisplayName:          directoryEntity.DisplayName,
			IdentityProviderType: &directoryEntity.IdentityProviderType,
			Type_:                directoryEntity.EntityType,
		})
	}

	for _, applicationId := range accessPolicy.Applications {
		applications = append(applications, sdk.ApplicationBase{
			Id:    applicationId,
			Type_: ToApplicationType(accessPolicy.TargetProtocol),
		})
	}

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

		if accessPolicy.Validators.ComplianceCheck {
			validatorsDto[ValidatorComplianceCheck] = accessPolicy.Validators.ComplianceCheck
		}

		if accessPolicy.Validators.WebVerification {
			validatorsDto[ValidatorWebVerification] = accessPolicy.Validators.WebVerification
		}
	}

	if accessPolicy.Conditions != nil {
		if accessPolicy.Conditions.SourceIp != nil {
			conditionsDto = append(conditionsDto, sdk.PolicyCondition{
				ConditionDefinitionId: IpCondition,
				Arguments: map[string][]string{
					IpUuid:           accessPolicy.Conditions.SourceIp,
					SharedIpListUuid: accessPolicy.Conditions.SharedIpList,
				},
			})
		}

		if accessPolicy.Conditions.Location != nil {
			conditionsDto = append(conditionsDto, sdk.PolicyCondition{
				ConditionDefinitionId: LocationRestrictionCondition,
				Arguments: map[string][]string{
					CountriesUuid: accessPolicy.Conditions.Location,
				},
			})
		}

		if accessPolicy.Conditions.ManagedDevice.SymantecCloudSoc || accessPolicy.Conditions.ManagedDevice.SymantecWebSecurityService {
			var managedDeviceArguments []string

			if accessPolicy.Conditions.ManagedDevice.SymantecWebSecurityService {
				managedDeviceArguments = append(managedDeviceArguments, ManagedDeviceWssConditionArgument)
			}

			if accessPolicy.Conditions.ManagedDevice.SymantecCloudSoc {
				managedDeviceArguments = append(managedDeviceArguments, ManagedDeviceCloudSocConditionArgument)
			}

			if accessPolicy.Conditions.ManagedDevice.OpswatMetaAccess {
				managedDeviceArguments = append(managedDeviceArguments, ManagedDeviceOpswatConditionArgument)
			}

			conditionsDto = append(conditionsDto, sdk.PolicyCondition{
				ConditionDefinitionId: ManagedDeviceCondition,
				Arguments:             map[string][]string{ManagedDeviceUuid: managedDeviceArguments},
			})
		}

		if accessPolicy.Conditions.UnmanagedDevice {
			conditionsDto = append(conditionsDto, sdk.PolicyCondition{
				ConditionDefinitionId: UnmanagedDeviceCondition,
			})
		}
	}

	accessPolicyDto := sdk.PolicyAccess{
		Type_:             &accessPolicyType,
		TargetProtocol:    ToTargetProtocol(accessPolicy.TargetProtocol),
		Id:                accessPolicy.Id,
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

	for _, directoryEntityDto := range accessPolicyDto.DirectoryEntities {
		directoryEntity = append(directoryEntity, DirectoryEntity{
			IdentifierInProvider: directoryEntityDto.IdentifierInProvider,
			IdentityProviderId:   directoryEntityDto.IdentityProviderId,
			DisplayName:          directoryEntityDto.DisplayName,
			IdentityProviderType: *directoryEntityDto.IdentityProviderType,
			EntityType:           directoryEntityDto.Type_,
		})
	}

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
			ComplianceCheck: accessPolicyDto.Validators[ValidatorComplianceCheck],
			WebVerification: accessPolicyDto.Validators[ValidatorWebVerification],
		}
	}

	if accessPolicyDto.FilterConditions != nil && len(accessPolicyDto.FilterConditions) > 0 {
		conditions = &Conditions{}

		for _, filterCondition := range accessPolicyDto.FilterConditions {
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

			if filterCondition.ConditionDefinitionId == ManagedDeviceCondition {
				for _, managedDeviceCondition := range filterCondition.Arguments[ManagedDeviceUuid] {

					if ManagedDeviceOpswatConditionArgument == managedDeviceCondition {
						conditions.ManagedDevice.OpswatMetaAccess = true
					}

					if ManagedDeviceCloudSocConditionArgument == managedDeviceCondition {
						conditions.ManagedDevice.SymantecCloudSoc = true
					}

					if ManagedDeviceWssConditionArgument == managedDeviceCondition {
						conditions.ManagedDevice.SymantecWebSecurityService = true
					}
				}
			}

			if filterCondition.ConditionDefinitionId == UnmanagedDeviceCondition {
				conditions.UnmanagedDevice = true
			}
		}
	}

	return &AccessPolicy{
		TargetProtocol:    FromTargetProtocol(*accessPolicyDto.TargetProtocol),
		Id:                accessPolicyDto.Id,
		Enabled:           accessPolicyDto.Enabled,
		CreatedAt:         accessPolicyDto.CreatedAt,
		Name:              accessPolicyDto.Name,
		DirectoryEntities: directoryEntity,
		Applications:      applications,
		Validators:        validators,
		Conditions:        conditions,
		RdpSettings:       rdpSetting,
		SshSettings:       sshSetting,
		TcpSettings:       tcpSetting,
	}
}
