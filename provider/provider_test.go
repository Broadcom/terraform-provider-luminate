package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var newTestAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider
var testAccDomain string

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"luminate": testAccProvider,
	}
	newTestAccProviders = make(map[string]func() (*schema.Provider, error))
	newTestAccProviders["luminate"] = func() (*schema.Provider, error) {
		return testAccProvider, nil
	}
	testAccDomain = "terraformat.luminatesec.com"
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("LUMINATE_API_ENDPOINT"); v == "" {
		t.Fatal("LUMINATE_API_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("LUMINATE_API_CLIENT_ID"); v == "" {
		t.Fatal("LUMINATE_API_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("LUMINATE_API_CLIENT_SECRET"); v == "" {
		t.Fatal("LUMINATE_API_CLIENT_SECRET must be set for acceptance tests")
	}
}
