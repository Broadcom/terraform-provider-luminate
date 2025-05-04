package provider

import (
	"os"
	"strings"
	"testing"

	"github.com/Broadcom/terraform-provider-luminate/framework_provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtocol6Providers map[string]func() (tfprotov6.ProviderServer, error)
var testAccDomain string

func init() {
	testAccProtocol6Providers = map[string]func() (tfprotov6.ProviderServer, error){
		"luminate": func() (tfprotov6.ProviderServer, error) {
			providerServer, err := framework_provider.CreateProviderServer(Provider())
			if err != nil {
				return nil, err
			}

			return providerServer(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	apiEndpoint := os.Getenv("LUMINATE_API_ENDPOINT")
	if apiEndpoint == "" {
		t.Fatal("LUMINATE_API_ENDPOINT must be set for acceptance tests")
	} else {
		testAccDomain = strings.Replace(apiEndpoint, "api.", "", 1)
	}
	if v := os.Getenv("LUMINATE_API_CLIENT_ID"); v == "" {
		t.Fatal("LUMINATE_API_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("LUMINATE_API_CLIENT_SECRET"); v == "" {
		t.Fatal("LUMINATE_API_CLIENT_SECRET must be set for acceptance tests")
	}
}
