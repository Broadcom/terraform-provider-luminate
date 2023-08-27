package provider

import (
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeviceList(t *testing.T) {

	testCases := []struct {
		input    []interface{}
		expected dto.Device
	}{
		// Define your test cases here
		{
			[]interface{}{
				map[string]interface{}{
					"opswat":                        true,
					"symantec_cloudsoc":             true,
					"symantec_web_security_service": true},
			},
			dto.Device{
				OpswatMetaAccess:           true,
				SymantecCloudSoc:           true,
				SymantecWebSecurityService: true,
			}, // Sum of ages: 25 + 30 = 55
		},
		{
			[]interface{}{
				map[string]interface{}{
					"opswat":                        true,
					"symantec_cloudsoc":             false,
					"symantec_web_security_service": false},
			},
			dto.Device{
				OpswatMetaAccess:           true,
				SymantecCloudSoc:           false,
				SymantecWebSecurityService: false,
			}, // Sum of ages: 25 + 30 = 55
		},
	}
	// Test each test case
	for _, tc := range testCases {
		var actual dto.Device
		deviceList(tc.input, &actual)
		assert.Equal(t, actual.SymantecWebSecurityService, tc.expected.SymantecWebSecurityService)
		assert.Equal(t, actual.SymantecCloudSoc, tc.expected.SymantecCloudSoc)
		assert.Equal(t, actual.OpswatMetaAccess, tc.expected.OpswatMetaAccess)
	}
}
