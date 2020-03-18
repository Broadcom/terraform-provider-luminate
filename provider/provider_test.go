package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccDomain string

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"luminate": testAccProvider,
	}
	testAccDomain = "tfacc.luminatesec.com"
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
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
