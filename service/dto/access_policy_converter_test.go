package dto

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	sdk "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"
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

func TestToTargetProtocolSubType(t *testing.T) {
	tests := []struct {
		name                          string
		targetProtocolSubTypeString   string
		expectedTargetProtocolSubType sdk.PolicyTargetProtocolSubType
	}{
		{"RDP_NATIVE", "RDP_NATIVE", sdk.NATIVE_PolicyTargetProtocolSubType},
		{"RDP_BROWSER", "RDP_BROWSER", sdk.BROWSER_PolicyTargetProtocolSubType},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedTargetProtocolSubType, *ToTargetProtocolSubType(test.targetProtocolSubTypeString))
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

func TestToApplicationSubType(t *testing.T) {
	tests := []struct {
		name                       string
		applicationSubTypeString   string
		expectedApplicationSubType sdk.ApplicationSubType
	}{
		{"HTTP_LUMINATE_DOMAIN", "HTTP_LUMINATE_DOMAIN", sdk.HTTP_LUMINATE_DOMAIN_ApplicationSubType},
		{"HTTP_CUSTOM_DOMAIN", "HTTP_CUSTOM_DOMAIN", sdk.HTTP_CUSTOM_DOMAIN_ApplicationSubType},
		{"HTTP_WILDCARD_DOMAIN", "HTTP_WILDCARD_DOMAIN", sdk.HTTP_WILDCARD_DOMAIN_ApplicationSubType},
		{"SINGLE_MACHINE", "SINGLE_MACHINE", sdk.SINGLE_MACHINE_ApplicationSubType},
		{"MULTIPLE_MACHINES", "MULTIPLE_MACHINES", sdk.MULTIPLE_MACHINES_ApplicationSubType},
		{"RDP_BROWSER_SINGLE_MACHINE", "RDP_BROWSER_SINGLE_MACHINE", sdk.RDP_BROWSER_SINGLE_MACHINE_ApplicationSubType},
		{"RDP_BROWSER_MULTIPLE_MACHINES", "RDP_BROWSER_MULTIPLE_MACHINES", sdk.RDP_BROWSER_MULTIPLE_MACHINES_ApplicationSubType},
		{"SEGMENT_SPECIFIC_IPS", "SEGMENT_SPECIFIC_IPS", sdk.SEGMENT_SPECIFIC_IPS_ApplicationSubType},
		{"SEGMENT_RANGE", "SEGMENT_RANGE", sdk.SEGMENT_RANGE_ApplicationSubType},
		{"Unknown", "Unknown", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedApplicationSubType, *ToApplicationSubType(test.applicationSubTypeString))
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

// mockAppFetcher is a mock implementation of the ApplicationFetcher interface for testing.
type mockAppFetcher struct {
	mockApp *Application
	mockErr error
}

func (m *mockAppFetcher) GetApplicationById(applicationID string) (*Application, error) {
	if m.mockErr != nil {
		return nil, m.mockErr
	}
	if m.mockApp == nil {
		return nil, nil
	}

	app := *m.mockApp
	app.ID = applicationID
	return &app, nil
}

func TestConvertToDto(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		appID := uuid.New().String()
		accessPolicy := &AccessPolicy{
			Policy: Policy{
				Enabled:        true,
				Name:           "my policy",
				Applications:   []string{appID},
				TargetProtocol: "SSH",
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

		// setup mock fetcher
		mockFetcher := &mockAppFetcher{
			mockApp: &Application{
				ID:      appID,
				SubType: string(sdk.SINGLE_MACHINE_ApplicationSubType),
			},
		}

		// when
		accessPolicyDto, err := ConvertToDto(accessPolicy, mockFetcher)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, accessPolicyDto)

		// Check that the application subtype was set correctly
		assert.Len(t, accessPolicyDto.Applications, 1)
		assert.Equal(t, sdk.SINGLE_MACHINE_ApplicationSubType, *accessPolicyDto.Applications[0].SubType)

		// Check that the conversion back and forth works
		generatedAccessPolicy := ConvertFromDto(*accessPolicyDto)
		for i := range generatedAccessPolicy.DirectoryEntities {
			generatedAccessPolicy.DirectoryEntities[i].IdentityProviderType = ""
		}

		assert.Equal(t, accessPolicy, generatedAccessPolicy)
	})

	t.Run("success with WebRDP custom settings", func(t *testing.T) {
		// given
		appID := uuid.New().String()
		accessPolicy := &AccessPolicy{
			Policy: Policy{
				Enabled:        true,
				Name:           "my policy",
				Applications:   []string{appID},
				TargetProtocol: "RDP",
			},
			RdpSettings: &PolicyRdpSettings{
				LongTermPassword: true,
				WebRdpSettings: &PolicyWebRdpSettings{
					DisableCopy:  true,
					DisablePaste: true,
				},
			},
		}

		// setup mock fetcher
		mockFetcher := &mockAppFetcher{
			mockApp: &Application{
				ID:      appID,
				SubType: string(sdk.SINGLE_MACHINE_ApplicationSubType),
			},
		}

		// when
		accessPolicyDto, err := ConvertToDto(accessPolicy, mockFetcher)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, accessPolicyDto)

		// Check that the application subtype was set correctly
		assert.Len(t, accessPolicyDto.Applications, 1)
		assert.Equal(t, sdk.SINGLE_MACHINE_ApplicationSubType, *accessPolicyDto.Applications[0].SubType)

		// Check that the conversion back and forth works
		generatedAccessPolicy := ConvertFromDto(*accessPolicyDto)
		for i := range generatedAccessPolicy.DirectoryEntities {
			generatedAccessPolicy.DirectoryEntities[i].IdentityProviderType = ""
		}

		assert.Equal(t, accessPolicy, generatedAccessPolicy)
	})

	t.Run("error when fetching application fails", func(t *testing.T) {
		// given
		accessPolicy := &AccessPolicy{Policy: Policy{Applications: []string{uuid.New().String()}}}
		expectedErr := errors.New("failed to fetch")
		mockFetcher := &mockAppFetcher{mockErr: expectedErr}

		// when
		dto, err := ConvertToDto(accessPolicy, mockFetcher)

		// then
		assert.Error(t, err)
		assert.Nil(t, dto)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})

	t.Run("error when application is not found", func(t *testing.T) {
		// given
		accessPolicy := &AccessPolicy{Policy: Policy{Applications: []string{uuid.New().String()}}}
		mockFetcher := &mockAppFetcher{mockApp: nil} // Simulates app not found

		// when
		dto, err := ConvertToDto(accessPolicy, mockFetcher)

		// then
		assert.Error(t, err)
		assert.Nil(t, dto)
		assert.Contains(t, err.Error(), "not found")
	})
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
