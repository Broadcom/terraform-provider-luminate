package framework_provider

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccResourceSiteRegistrationKeyVersionTemplate = `
resource "luminate_site_registration_key_version" "site_registration_key_version" {
  version = <VERSION_PLACEHOLDER>
}
`

const testAccResourceSiteRegistrationKeyVersionTemplateWithSiteToForceApply = `
resource "luminate_site" "new_site_<RANDOM_PLACEHOLDER>" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
  	authentication_mode = "orchestrator"
}

resource "luminate_site_registration_key_version" "site_registration_key_version" {
  version = <VERSION_PLACEHOLDER>
}
`

func TestAccLuminateSiteRegistrationKeyVersion(t *testing.T) {
	firstVersion := strings.ReplaceAll(testAccResourceSiteRegistrationKeyVersionTemplate, "<VERSION_PLACEHOLDER>", "1")
	secondVersion := strings.ReplaceAll(testAccResourceSiteRegistrationKeyVersionTemplate, "<VERSION_PLACEHOLDER>", "2")

	secondVersionWithSiteToForceApply := strings.ReplaceAll(testAccResourceSiteRegistrationKeyVersionTemplateWithSiteToForceApply, "<VERSION_PLACEHOLDER>", "2")
	secondVersionWithSiteToForceApply = strings.ReplaceAll(secondVersionWithSiteToForceApply, "<RANDOM_PLACEHOLDER>", strconv.Itoa(100+rand.Intn(100)))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: firstVersion,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("luminate_site_registration_key_version.site_registration_key_version", "version", "1"),
					resource.TestCheckResourceAttr("luminate_site_registration_key_version.site_registration_key_version", "version_changed", "true"),
				),
			},
			{
				Config:   firstVersion,
				PlanOnly: true,
			},
			{
				Config: secondVersion,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("luminate_site_registration_key_version.site_registration_key_version", "version", "2"),
					resource.TestCheckResourceAttr("luminate_site_registration_key_version.site_registration_key_version", "version_changed", "true"),
				),
			},
			{
				Config: secondVersionWithSiteToForceApply,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("luminate_site_registration_key_version.site_registration_key_version", "version", "2"),
					resource.TestCheckResourceAttr("luminate_site_registration_key_version.site_registration_key_version", "version_changed", "false"),
				),
			},
		},
	})
}
