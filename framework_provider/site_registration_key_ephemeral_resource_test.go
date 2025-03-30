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

ephemeral "luminate_site_registration_key" "new_site_registration_key" {
	site_id = luminate_site.new_site_<RANDOM_PLACEHOLDER>.id
	revoke_existing_key_immediately = true
	rotate = true
}
`

func TestAccLuminateSiteRegistrationKey(t *testing.T) {
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: strings.ReplaceAll(testAccEphemeralResourceSiteRegistrationKeyTemplate, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
			},
		},
	})
}
