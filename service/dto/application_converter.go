package dto

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"log"
)

func ConvertToApplicationDTO(applicationSDKDTO sdk.Application) Application {
	applicationServiceDTO := Application{
		ID:                   applicationSDKDTO.Id,
		Name:                 applicationSDKDTO.Name,
		CollectionID:         applicationSDKDTO.CollectionId,
		Visible:              applicationSDKDTO.IsVisible,
		NotificationsEnabled: applicationSDKDTO.IsNotificationEnabled,
		Type:                 GetApplicationTypeString(*applicationSDKDTO.Type_),
	}

	if applicationSDKDTO.SubType != nil {
		applicationServiceDTO.SubType = string(*applicationSDKDTO.SubType)
	}

	if applicationSDKDTO.ConnectionSettings != nil {
		applicationServiceDTO.Subdomain = applicationSDKDTO.ConnectionSettings.Subdomain
		applicationServiceDTO.InternalAddress = applicationSDKDTO.ConnectionSettings.InternalAddress
		applicationServiceDTO.ExternalAddress = applicationSDKDTO.ConnectionSettings.ExternalAddress
		applicationServiceDTO.LuminateAddress = applicationSDKDTO.ConnectionSettings.LuminateAddress
		applicationServiceDTO.CustomRootPath = applicationSDKDTO.ConnectionSettings.CustomRootPath
		applicationServiceDTO.CustomExternalAddress = applicationSDKDTO.ConnectionSettings.CustomExternalAddress
		applicationServiceDTO.WildcardCertificate = applicationSDKDTO.ConnectionSettings.CustomSSLCertificate
		applicationServiceDTO.WildcardPrivateKey = applicationSDKDTO.ConnectionSettings.WildcardPrivateKey
	}

	if applicationSDKDTO.LinkTranslationSettings != nil {
		applicationServiceDTO.DefaultContentRewriteRulesEnabled = applicationSDKDTO.LinkTranslationSettings.IsDefaultContentRewriteRulesEnabled
		applicationServiceDTO.DefaultHeaderRewriteRulesEnabled = applicationSDKDTO.LinkTranslationSettings.IsDefaultHeaderRewriteRulesEnabled
		applicationServiceDTO.UseExternalAddressForHostAndSni = applicationSDKDTO.LinkTranslationSettings.UseExternalAddressForHostAndSni

		var linkedins []*string
		for index := range applicationSDKDTO.LinkTranslationSettings.LinkedApplications {
			linkedins = append(linkedins, &applicationSDKDTO.LinkTranslationSettings.LinkedApplications[index])
		}
		applicationServiceDTO.LinkedApplications = linkedins
	}
	if applicationSDKDTO.RequestCustomizationSettings != nil {
		applicationServiceDTO.HeaderCustomization = HeaderStringsToMap(*applicationSDKDTO.RequestCustomizationSettings.HeaderCustomization)
	}

	tcpTunnelSettings := applicationSDKDTO.TcpTunnelSettings
	if applicationSDKDTO.TcpTunnelSettings != nil && len(*tcpTunnelSettings) > 0 {
		for _, t := range *tcpTunnelSettings {
			target := TCPTarget{
				Address:     t.Target,
				Ports:       t.Ports,
				PortMapping: t.PortMapping,
			}
			applicationServiceDTO.Targets = append(applicationServiceDTO.Targets, target)
		}
	}

	if applicationSDKDTO.CloudIntegrationData != nil {
		tags := map[string]string{}
		var vpcs []Vpc

		for _, tagSdk := range applicationSDKDTO.CloudIntegrationData.Tags {
			tags[tagSdk.Key] = tagSdk.Value
		}

		for _, vpcSdk := range applicationSDKDTO.CloudIntegrationData.Vpcs {
			vpcs = append(vpcs, Vpc{
				IntegrationId: vpcSdk.IntegrationId,
				Region:        vpcSdk.Region,
				Vpc:           vpcSdk.Vpc,
				CidrBlock:     vpcSdk.CidrBlock,
			})
		}

		cloudIntegrationData := &CloudIntegrationData{
			SegmentId: applicationSDKDTO.CloudIntegrationData.SegmentId,
			Tags:      tags,
			Vpcs:      vpcs,
		}

		applicationServiceDTO.CloudIntegrationData = cloudIntegrationData
	}

	if applicationSDKDTO.SegmentSettings != nil {
		applicationServiceDTO.SegmentSettings = &SegmentSettings{OriginalIP: applicationSDKDTO.SegmentSettings.OriginalIp}
	}
	if applicationSDKDTO.MultipleSegmentSettings != nil {
		var multipleSegmentSettings []*SegmentSettings
		for i := 0; i < len(*applicationSDKDTO.MultipleSegmentSettings); i++ {
			multipleSegmentSettings = append(multipleSegmentSettings,
				&SegmentSettings{OriginalIP: (*applicationSDKDTO.MultipleSegmentSettings)[i].OriginalIp})
		}
		applicationServiceDTO.MultipleSegmentSettings = multipleSegmentSettings
	}

	return applicationServiceDTO
}

func ConvertFromApplicationDTO(applicationServiceDTO Application) sdk.Application {

	aType := GetApplicationType(applicationServiceDTO.Type)

	applicationSDKDTO := sdk.Application{
		Name:                  applicationServiceDTO.Name,
		CollectionId:          applicationServiceDTO.CollectionID,
		Type_:                 &aType,
		Icon:                  applicationServiceDTO.Icon,
		Enabled:               true,
		IsVisible:             applicationServiceDTO.Visible,
		IsNotificationEnabled: applicationServiceDTO.NotificationsEnabled,
		ConnectionSettings: &sdk.ApplicationConnectionSettings{
			InternalAddress:      applicationServiceDTO.InternalAddress,
			ExternalAddress:      applicationServiceDTO.ExternalAddress,
			CustomRootPath:       applicationServiceDTO.CustomRootPath,
			Subdomain:            applicationServiceDTO.Subdomain,
			LuminateAddress:      applicationServiceDTO.LuminateAddress,
			CustomSSLCertificate: applicationServiceDTO.WildcardCertificate,
			WildcardPrivateKey:   applicationServiceDTO.WildcardPrivateKey,
		},
	}
	if applicationServiceDTO.SubType != "" {
		saType := GetApplicationSubType(applicationServiceDTO.SubType)
		applicationSDKDTO.SubType = &saType
	}

	switch applicationServiceDTO.Type {
	case "web":
		var linkedApps []string
		for _, v := range applicationServiceDTO.LinkedApplications {
			linkedApps = append(linkedApps, *v)
		}

		aSubType := GetApplicationSubType(applicationServiceDTO.SubType)
		applicationSDKDTO.SubType = &aSubType

		if applicationServiceDTO.CustomExternalAddress != "" {
			applicationSDKDTO.ConnectionSettings.CustomExternalAddress = applicationServiceDTO.CustomExternalAddress
		}

		applicationSDKDTO.LinkTranslationSettings = &sdk.ApplicationLinkTranslationSettings{
			IsDefaultContentRewriteRulesEnabled: applicationServiceDTO.DefaultContentRewriteRulesEnabled,
			IsDefaultHeaderRewriteRulesEnabled:  applicationServiceDTO.DefaultHeaderRewriteRulesEnabled,
			UseExternalAddressForHostAndSni:     applicationServiceDTO.UseExternalAddressForHostAndSni,
			LinkedApplications:                  linkedApps,
		}

		headers := HeaderMapToStrings(applicationServiceDTO.HeaderCustomization)
		applicationSDKDTO.RequestCustomizationSettings = &sdk.ApplicationRequestCustomizationSettings{
			HeaderCustomization: &headers,
		}
	case "tcp":
		ApplicationTcpTargetSlice := make([]sdk.ApplicationTcpTarget, 0)
		for _, v := range applicationServiceDTO.Targets {
			t := sdk.ApplicationTcpTarget{
				Ports:       v.Ports,
				Target:      v.Address,
				PortMapping: v.PortMapping,
			}
			log.Printf("[DEBUG] TUNNEL Target %v", t)

			ApplicationTcpTargetSlice = append(ApplicationTcpTargetSlice, t)
		}
		applicationSDKDTO.TcpTunnelSettings = &ApplicationTcpTargetSlice

		log.Printf("[DEBUG] TUNNEL SETTINGS %v", applicationSDKDTO.TcpTunnelSettings)
	case "sshgw":
		var sdkVpcs []sdk.ApplicationVpcData
		for _, vpc := range applicationServiceDTO.CloudIntegrationData.Vpcs {
			sdkVpcs = append(sdkVpcs, sdk.ApplicationVpcData{
				IntegrationId: vpc.IntegrationId,
				Vpc:           vpc.Vpc,
				Region:        vpc.Region,
				CidrBlock:     vpc.CidrBlock,
			})
		}

		var tagsSdk []sdk.ApplicationCloudIntegrationTag
		for key, value := range applicationServiceDTO.CloudIntegrationData.Tags {
			tagsSdk = append(tagsSdk, sdk.ApplicationCloudIntegrationTag{
				Key:   key,
				Value: value,
			})
		}

		applicationSDKDTO.CloudIntegrationData = &sdk.ApplicationCloudIntegrationDataProperties{
			Tags:      tagsSdk,
			Vpcs:      sdkVpcs,
			SegmentId: applicationServiceDTO.CloudIntegrationData.SegmentId,
		}
	case "segment":
		if applicationServiceDTO.SegmentSettings != nil {
			applicationSDKDTO.SegmentSettings = &sdk.ApplicationConnectionSettingsSegment{
				OriginalIp: applicationServiceDTO.SegmentSettings.OriginalIP,
			}
		}
		if applicationServiceDTO.MultipleSegmentSettings != nil {
			var multipleSegmentSettings []sdk.ApplicationConnectionSettingsSegment
			for i := 0; i < len(applicationServiceDTO.MultipleSegmentSettings); i++ {
				multipleSegmentSettings = append(multipleSegmentSettings, sdk.ApplicationConnectionSettingsSegment{
					OriginalIp: applicationServiceDTO.MultipleSegmentSettings[i].OriginalIP,
				})
			}
			applicationSDKDTO.MultipleSegmentSettings = &multipleSegmentSettings
		}
	}

	return applicationSDKDTO
}

func HeaderMapToStrings(headers map[string]interface{}) []map[string]string {
	var result []map[string]string

	for k, v := range headers {
		result = append(result, map[string]string{k: v.(string)})
	}

	return result
}

func HeaderStringsToMap(headers []map[string]string) map[string]interface{} {
	var result map[string]interface{}

	for i, header := range headers {
		if i == 0 {
			result = make(map[string]interface{})
		}
		for k, v := range header {
			result[k] = v
		}

	}

	return result
}

func GetApplicationType(appType string) sdk.ApplicationType {
	switch appType {
	case "web":
		return sdk.HTTP_ApplicationType
	case "ssh":
		return sdk.SSH_ApplicationType
	case "sshgw":
		return sdk.DYNAMIC_SSH_ApplicationType
	case "tcp":
		return sdk.TCP_ApplicationType
	case "rdp":
		return sdk.RDP_ApplicationType
	case "segment":
		return sdk.SEGMENT_ApplicationType
	}
	return ""
}

func GetApplicationTypeString(appType sdk.ApplicationType) string {
	switch appType {
	case sdk.HTTP_ApplicationType:
		return "web"
	case sdk.SSH_ApplicationType:
		return "ssh"
	case sdk.DYNAMIC_SSH_ApplicationType:
		return "sshgw"
	case sdk.TCP_ApplicationType:
		return "tcp"
	case sdk.RDP_ApplicationType:
		return "rdp"
	case sdk.SEGMENT_ApplicationType:
		return "segment"
	}
	return ""
}

func GetApplicationSubType(appSubType string) sdk.ApplicationSubType {
	switch appSubType {
	case string(sdk.HTTP_LUMINATE_DOMAIN_ApplicationSubType):
		return sdk.HTTP_LUMINATE_DOMAIN_ApplicationSubType
	case string(sdk.HTTP_CUSTOM_DOMAIN_ApplicationSubType):
		return sdk.HTTP_CUSTOM_DOMAIN_ApplicationSubType
	case string(sdk.HTTP_WILDCARD_DOMAIN_ApplicationSubType):
		return sdk.HTTP_WILDCARD_DOMAIN_ApplicationSubType
	case string(sdk.SINGLE_MACHINE_ApplicationSubType):
		return sdk.SINGLE_MACHINE_ApplicationSubType
	case string(sdk.MULTIPLE_MACHINES_ApplicationSubType):
		return sdk.MULTIPLE_MACHINES_ApplicationSubType
	case string(sdk.SEGMENT_RANGE_ApplicationSubType):
		return sdk.SEGMENT_RANGE_ApplicationSubType
	case string(sdk.SEGMENT_SPECIFIC_IPS_ApplicationSubType):
		return sdk.SEGMENT_SPECIFIC_IPS_ApplicationSubType
	}
	return ""
}
