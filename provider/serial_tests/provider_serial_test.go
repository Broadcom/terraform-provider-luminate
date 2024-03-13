package serial_tests

import (
	"github.com/Broadcom/terraform-provider-luminate/provider"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var newTestAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = provider.Provider()
	newTestAccProviders = make(map[string]func() (*schema.Provider, error))
	newTestAccProviders["luminate"] = func() (*schema.Provider, error) {
		return testAccProvider, nil
	}
}

func TestProvider(t *testing.T) {
	if err := provider.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	apiEndpoint := os.Getenv("LUMINATE_API_ENDPOINT")
	if apiEndpoint == "" {
		t.Fatal("LUMINATE_API_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("LUMINATE_API_CLIENT_ID"); v == "" {
		t.Fatal("LUMINATE_API_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("LUMINATE_API_CLIENT_SECRET"); v == "" {
		t.Fatal("LUMINATE_API_CLIENT_SECRET must be set for acceptance tests")
	}
}
