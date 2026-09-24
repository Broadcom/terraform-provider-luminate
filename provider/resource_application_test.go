package provider

import (
	"testing"

	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSetSSHApplicationFieldsNormalizesInternalAddress(t *testing.T) {
	testCases := []struct {
		name       string
		fromServer string
		expected   string
	}{
		{
			name:       "default port added by the server is not persisted",
			fromServer: "tcp://127.0.0.1:22",
			expected:   "tcp://127.0.0.1",
		},
		{
			name:       "address returned without a port is persisted as is",
			fromServer: "tcp://127.0.0.1",
			expected:   "tcp://127.0.0.1",
		},
		{
			name:       "explicit non default port is persisted",
			fromServer: "tcp://127.0.0.1:2222",
			expected:   "tcp://127.0.0.1:2222",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, LuminateSSHApplication().Schema, map[string]interface{}{})

			setSSHApplicationFields(d, &dto.Application{InternalAddress: tc.fromServer})

			if got := d.Get("internal_address").(string); got != tc.expected {
				t.Errorf("expected internal_address %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestStripDefaultPort(t *testing.T) {
	testCases := []struct {
		name        string
		address     string
		defaultPort string
		expected    string
	}{
		{
			name:        "ssh default port is removed",
			address:     "tcp://127.0.0.1:22",
			defaultPort: utils.DefaultSSHPort,
			expected:    "tcp://127.0.0.1",
		},
		{
			name:        "explicit non default port is kept",
			address:     "tcp://127.0.0.1:2222",
			defaultPort: utils.DefaultSSHPort,
			expected:    "tcp://127.0.0.1:2222",
		},
		{
			name:        "port with the default as a prefix is kept",
			address:     "tcp://127.0.0.1:220",
			defaultPort: utils.DefaultSSHPort,
			expected:    "tcp://127.0.0.1:220",
		},
		{
			name:        "address without a port is unchanged",
			address:     "tcp://127.0.0.1",
			defaultPort: utils.DefaultSSHPort,
			expected:    "tcp://127.0.0.1",
		},
		{
			name:        "another type's default port is kept",
			address:     "tcp://127.0.0.1:3389",
			defaultPort: utils.DefaultSSHPort,
			expected:    "tcp://127.0.0.1:3389",
		},
		{
			name:        "rdp default port is removed",
			address:     "tcp://127.0.0.1:3389",
			defaultPort: utils.DefaultRDPPort,
			expected:    "tcp://127.0.0.1",
		},
		{
			name:        "address without a scheme is unchanged",
			address:     "127.0.0.1",
			defaultPort: utils.DefaultSSHPort,
			expected:    "127.0.0.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := stripDefaultPort(tc.address, tc.defaultPort); got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestSuppressDefaultPortDiff(t *testing.T) {
	testCases := []struct {
		name        string
		defaultPort string
		oldValue    string
		newValue    string
		suppressed  bool
	}{
		{
			name:        "empty old value is a real change",
			defaultPort: utils.DefaultSSHPort,
			oldValue:    "",
			newValue:    "tcp://127.0.0.1",
			suppressed:  false,
		},
		{
			name:        "identical values",
			defaultPort: utils.DefaultSSHPort,
			oldValue:    "tcp://127.0.0.1:22",
			newValue:    "tcp://127.0.0.1:22",
			suppressed:  true,
		},
		{
			name:        "server appended the default port",
			defaultPort: utils.DefaultSSHPort,
			oldValue:    "tcp://127.0.0.1:22",
			newValue:    "tcp://127.0.0.1",
			suppressed:  true,
		},
		{
			name:        "configuration states the default port explicitly",
			defaultPort: utils.DefaultSSHPort,
			oldValue:    "tcp://127.0.0.1",
			newValue:    "tcp://127.0.0.1:22",
			suppressed:  true,
		},
		{
			name:        "different host is a real change",
			defaultPort: utils.DefaultSSHPort,
			oldValue:    "tcp://127.0.0.1",
			newValue:    "tcp://127.0.0.2",
			suppressed:  false,
		},
		{
			name:        "different explicit ports is a real change",
			defaultPort: utils.DefaultSSHPort,
			oldValue:    "tcp://127.0.0.1:2222",
			newValue:    "tcp://127.0.0.1:2223",
			suppressed:  false,
		},
		{
			name:        "dropping a non default port is a real change",
			defaultPort: utils.DefaultSSHPort,
			oldValue:    "tcp://127.0.0.1:2222",
			newValue:    "tcp://127.0.0.1",
			suppressed:  false,
		},
		{
			name:        "rdp default port is suppressed for rdp",
			defaultPort: utils.DefaultRDPPort,
			oldValue:    "tcp://127.0.0.1:3389",
			newValue:    "tcp://127.0.0.1",
			suppressed:  true,
		},
		{
			name:        "rdp default port is a real change for ssh",
			defaultPort: utils.DefaultSSHPort,
			oldValue:    "tcp://127.0.0.1:3389",
			newValue:    "tcp://127.0.0.1",
			suppressed:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			suppress := suppressDefaultPortDiff(tc.defaultPort)
			if got := suppress("internal_address", tc.oldValue, tc.newValue, nil); got != tc.suppressed {
				t.Errorf("expected suppressed=%v, got %v", tc.suppressed, got)
			}
		})
	}
}
