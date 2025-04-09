package framework_provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

const testAccEphemeralResourceSiteRegistrationKeyTemplate = `
resource "luminate_site" "new_site_<RANDOM_PLACEHOLDER>" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
  	authentication_mode = "site"
}

resource "luminate_site_registration_key_version" "site_registration_key_version" {
}

ephemeral "luminate_site_registration_key" "new_site_registration_key" {
	site_id = luminate_site.new_site_<RANDOM_PLACEHOLDER>.id
	version = luminate_site_registration_key_version.site_registration_key_version.version
	revoke_existing_key_immediately = true
}
`

func TestAccLuminateSiteRegistrationKey(t *testing.T) {
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config:             strings.ReplaceAll(testAccEphemeralResourceSiteRegistrationKeyTemplate, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
