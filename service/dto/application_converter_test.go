package dto

import (
	"github.com/stretchr/testify/assert"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
	"testing"
)

func TestGetApplicationTypeWeb(t *testing.T) {
	assert.Equal(t, sdk.SSH_ApplicationType, GetApplicationType("ssh"))
}

func TestGetApplicationTypeSSH(t *testing.T) {
	assert.Equal(t, sdk.HTTP_ApplicationType, GetApplicationType("web"))

}

func TestGetApplicationTypeRDP(t *testing.T) {
	assert.Equal(t, sdk.RDP_ApplicationType, GetApplicationType("rdp"))

}

func TestGetApplicationTypeSSHGW(t *testing.T) {
	assert.Equal(t, sdk.DYNAMIC_SSH_ApplicationType, GetApplicationType("sshgw"))

}

func TestGetApplicationTypeTCP(t *testing.T) {
	assert.Equal(t, sdk.TCP_ApplicationType, GetApplicationType("tcp"))
}

func TestGetApplicationTypeBad(t *testing.T) {
	assert.Equal(t, sdk.ApplicationType(""), GetApplicationType("asd"))
}

func TestConvertApplicationDTO_Web(t *testing.T) {
	expected := Application{
		Name:                              "dummy name",
		Visible:                           true,
		Type:                              "web",
		NotificationsEnabled:              true,
		Subdomain:                         "subdomain",
		InternalAddress:                   "internal",
		ExternalAddress:                   "external",
		LuminateAddress:                   "lumaddr",
		CustomRootPath:                    "root",
		CustomExternalAddress:             "custom_external",
		DefaultHeaderRewriteRulesEnabled:  true,
		DefaultContentRewriteRulesEnabled: true,
		UseExternalAddressForHostAndSni:   true,
	}

	sdkDTO := ConvertFromApplicationDTO(expected)
	providerDTO := ConvertToApplicationDTO(sdkDTO)
	assert.Equal(t, expected, providerDTO)
}

func TestConvertApplicationDTO_SSH(t *testing.T) {
	expected := Application{
		Name:                 "dummy name",
		Visible:              true,
		Type:                 "ssh",
		NotificationsEnabled: true,
		Subdomain:            "subdomain",
		InternalAddress:      "internal",
		ExternalAddress:      "external",
		LuminateAddress:      "lumaddr",
	}

	sdkDTO := ConvertFromApplicationDTO(expected)
	providerDTO := ConvertToApplicationDTO(sdkDTO)
	assert.Equal(t, expected, providerDTO)
}

func TestConvertApplicationDTO_SshGw(t *testing.T) {
	// given
	tags := map[string]string{}
	tags["key"] = "value"

	expected := Application{
		Name:                 "dummy name",
		Visible:              true,
		Type:                 "sshgw",
		NotificationsEnabled: true,
		Subdomain:            "subdomain",
		InternalAddress:      "internal",
		ExternalAddress:      "external",
		LuminateAddress:      "lumaddr",
		CloudIntegrationData: &CloudIntegrationData{
			Tags:      tags,
			SegmentId: "segment-id",
			Vpcs: []Vpc{{
				IntegrationId: "integration_id",
				Region:        "region",
				Vpc:           "vpc-id",
				CidrBlock:     "1.1.1.1/18",
			},
				{
					IntegrationId: "integration_id-2",
					Region:        "region-2",
					Vpc:           "vpc-id-2",
					CidrBlock:     "1.1.1.1/19",
				},
			},
		},
	}

	// when
	sdkDTO := ConvertFromApplicationDTO(expected)
	providerDTO := ConvertToApplicationDTO(sdkDTO)

	// than
	assert.Equal(t, expected, providerDTO)
}

func TestConvertApplicationDTO_TCP(t *testing.T) {

	targets := []TCPTarget{{Address: "address", Ports: []int32{80}}}

	expected := Application{
		Name:                 "dummy name",
		Visible:              true,
		Type:                 "tcp",
		NotificationsEnabled: true,
		Subdomain:            "subdomain",
		InternalAddress:      "internal",
		ExternalAddress:      "external",
		LuminateAddress:      "lumaddr",
		Targets:              targets,
	}

	sdkDTO := ConvertFromApplicationDTO(expected)
	providerDTO := ConvertToApplicationDTO(sdkDTO)
	assert.Equal(t, expected, providerDTO)
}
