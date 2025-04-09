package framework_provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

const testAccResourceSiteRegistrationKeyVersionTemplate = `
resource "luminate_site_registration_key_version" "site_registration_key_version" {
}
`

func TestAccLuminateSiteRegistrationKeyVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config:             testAccResourceSiteRegistrationKeyVersionTemplate,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
