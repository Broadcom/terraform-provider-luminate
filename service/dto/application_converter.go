package dto

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"fmt"
	"log"
	"strings"
)

func ConvertToApplicationDTO(applicationSDKDTO sdk.Application) Application {
	applicationServiceDTO := Application{
		ID:                   applicationSDKDTO.Id,
		Name:                 applicationSDKDTO.Name,
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
		for _, v := range applicationSDKDTO.LinkTranslationSettings.LinkedApplications {
			linkedins = append(linkedins, &v)
		}
		applicationServiceDTO.LinkedApplications = linkedins
	}
	if applicationSDKDTO.RequestCustomizationSettings != nil {
		applicationServiceDTO.HeaderCustomization = HeaderStringsToMap(applicationSDKDTO.RequestCustomizationSettings.HeaderCustomization)
	}

	if len(applicationSDKDTO.TcpTunnelSettings) > 0 {
		for _, t := range applicationSDKDTO.TcpTunnelSettings {
			target := TCPTarget{
				Address: t.Target,
				Ports:   t.Ports,
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

	return applicationServiceDTO
}

func ConvertFromApplicationDTO(applicationServiceDTO Application) sdk.Application {
	aType := GetApplicationType(applicationServiceDTO.Type)

	applicationSDKDTO := sdk.Application{
		Name:                  applicationServiceDTO.Name,
		Type_:                 &aType,
		Icon:                  applicationServiceDTO.Icon,
		IsVisible:             applicationServiceDTO.Visible,
		IsNotificationEnabled: applicationServiceDTO.NotificationsEnabled,
		ConnectionSettings: &sdk.ApplicationConnectionSettings{
			InternalAddress:       applicationServiceDTO.InternalAddress,
			ExternalAddress:       applicationServiceDTO.ExternalAddress,
			CustomRootPath:        applicationServiceDTO.CustomRootPath,
			Subdomain:             applicationServiceDTO.Subdomain,
			LuminateAddress:       applicationServiceDTO.LuminateAddress,
			CustomSSLCertificate:  applicationServiceDTO.WildcardCertificate,
			WildcardPrivateKey:    applicationServiceDTO.WildcardPrivateKey,
		},
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

		applicationSDKDTO.RequestCustomizationSettings = &sdk.ApplicationRequestCustomizationSettings{
			HeaderCustomization: HeaderMapToStrings(applicationServiceDTO.HeaderCustomization),
		}
	case "tcp":
		applicationSDKDTO.TcpTunnelSettings = []sdk.ApplicationTcpTarget{}
		for _, v := range applicationServiceDTO.Targets {
			t := sdk.ApplicationTcpTarget{
				Ports:  v.Ports,
				Target: v.Address,
			}
			log.Printf("[DEBUG] TUNNEL Target %v", t)

			applicationSDKDTO.TcpTunnelSettings = append(applicationSDKDTO.TcpTunnelSettings, t)
		}
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
	}

	return applicationSDKDTO
}

func HeaderMapToStrings(headers map[string]interface{}) []string {
	var result []string
	for k, v := range headers {
		headerString := fmt.Sprintf("%s: %s", k, v)
		result = append(result, headerString)
	}
	return result
}

func HeaderStringsToMap(headers []string) map[string]interface{} {
	var result map[string]interface{}

	for _, v := range headers {
		headerMap := strings.Split(v, ":")
		if len(headerMap) == 2 {
			result[headerMap[0]] = headerMap[1]
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
	}
	return ""
}

func GetApplicationSubType(appSubType string) sdk.ApplicationSubType {
	switch appSubType {
	case string(sdk.LUMINATE_DOMAIN_ApplicationSubType):
		return sdk.LUMINATE_DOMAIN_ApplicationSubType
	case string(sdk.CUSTOM_DOMAIN_ApplicationSubType):
		return sdk.CUSTOM_DOMAIN_ApplicationSubType
	case string(sdk.WILDCARD_DOMAIN_ApplicationSubType):
		return sdk.WILDCARD_DOMAIN_ApplicationSubType
	}
	return ""
}