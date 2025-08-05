package dto

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
	"testing"
)

func TestFromTargetProtocol(t *testing.T) {
	tests := []struct {
		name                         string
		targetProtocol               sdk.PolicyTargetProtocol
		expectedTargetProtocolString string
	}{
		{"HTTP", sdk.HTTP_PolicyTargetProtocol, "HTTP"},
		{"RDP", sdk.RDP_PolicyTargetProtocol, "RDP"},
		{"SSH", sdk.SSH_PolicyTargetProtocol, "SSH"},
		{"TCP", sdk.TCP_PolicyTargetProtocol, "TCP"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedTargetProtocolString, FromTargetProtocol(test.targetProtocol))
		})
	}
}

func TestToTargetProtocol(t *testing.T) {
	tests := []struct {
		name                   string
		targetProtocolString   string
		expectedTargetProtocol sdk.PolicyTargetProtocol
	}{
		{"HTTP", "HTTP", sdk.HTTP_PolicyTargetProtocol},
		{"RDP", "RDP", sdk.RDP_PolicyTargetProtocol},
		{"SSH", "SSH", sdk.SSH_PolicyTargetProtocol},
		{"TCP", "TCP", sdk.TCP_PolicyTargetProtocol},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedTargetProtocol, *ToTargetProtocol(test.targetProtocolString))
		})
	}
}

func TestToApplicationType(t *testing.T) {
	tests := []struct {
		name                    string
		targetProtocolString    string
		expectedApplicationType sdk.ApplicationType
	}{
		{"HTTP", "HTTP", sdk.HTTP_ApplicationType},
		{"RDP", "RDP", sdk.RDP_ApplicationType},
		{"SSH", "SSH", sdk.SSH_ApplicationType},
		{"TCP", "TCP", sdk.TCP_ApplicationType},
		{"Unknown", "Unknown", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedApplicationType, *ToApplicationType(test.targetProtocolString))
		})
	}
}

func TestFromModalType(t *testing.T) {
	tests := []struct {
		name               string
		entityType         sdk.EntityType
		expectedEntityType string
	}{
		{"ApiClient", sdk.API_CLIENT_EntityType, "ApiClient"},
		{"Group", sdk.GROUP_EntityType, "Group"},
		{"User", sdk.USER_EntityType, "User"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedEntityType, FromModelType(test.entityType))
		})
	}
}

func TestToModalType(t *testing.T) {
	tests := []struct {
		name              string
		expectedModalType sdk.EntityType
		entityType        string
	}{
		{"ApiClient", sdk.API_CLIENT_EntityType, "ApiClient"},
		{"Group", sdk.GROUP_EntityType, "Group"},
		{"User", sdk.USER_EntityType, "User"},
		{"Unknown", "", "Unknown"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedModalType, *ToModelType(test.entityType))
		})
	}
}

func TestConvertToDto(t *testing.T) {
	// given
	accessPolicy := &AccessPolicy{
		Policy: Policy{
			Enabled:      true,
			Name:         "my policy",
			Applications: []string{uuid.New().String()},
		},

		Validators: &Validators{
			WebVerification: true,
		},

		SshSettings: &PolicySshSettings{
			Accounts:             []string{"ubuntu"},
			AutoMapping:          true,
			AgentForward:         false,
			AcceptTemporaryToken: false,
			AcceptCertificate:    true,
		},

		RdpSettings: &PolicyRdpSettings{
			LongTermPassword: true,
		},

		TcpSettings: &PolicyTcpSettings{
			AcceptTemporaryToken: true,
			AcceptCertificate:    true,
		},
	}

	// when
	accessPolicyDto := ConvertToDto(accessPolicy)

	// then
	generatedAccessPolicy := ConvertFromDto(accessPolicyDto)
	for i := range generatedAccessPolicy.DirectoryEntities {
		generatedAccessPolicy.DirectoryEntities[i].IdentityProviderType = ""
	}

	assert.Equal(t, accessPolicy, generatedAccessPolicy)
}

func TestGetDeviceArguments(t *testing.T) {
	device := Device{
		OpswatMetaAccess:           true,
		SymantecCloudSoc:           true,
		SymantecWebSecurityService: true,
	}

	expectedMap := map[string][]string{}
	expectedMap[Authentication] = []string{
		ManagedDeviceWssConditionArgument,
		ManagedDeviceCloudSocConditionArgument,
		ManagedDeviceOpswatConditionArgument,
	}

	actualMap := getDeviceArguments(device)

	assert.Equal(t, expectedMap, actualMap)
}

func TestHaveDeviceArgument(t *testing.T) {
	device := Device{
		OpswatMetaAccess:           false,
		SymantecCloudSoc:           false,
		SymantecWebSecurityService: false,
	}
	actual := hasDeviceArgument(device)
	assert.Equal(t, false, actual)

	device.SymantecWebSecurityService = true
	actual = hasDeviceArgument(device)

	assert.Equal(t, true, actual)

}
