package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testSiteRole = `
	resource "luminate_site" "site" {
		name = "siteToBeAssign"
	} 
	resource "luminate_site_role" "site-role" {
		role_type = "SiteEditor"
		identity_provider_id =  "local"
		entity_id = "a8a48219-835f-4183-a2a9-bbba8cad8eb8"
		entity_type = "User"
		site_id = "${luminate_site.site.id}"
	}
`

func TestAccLuminateSiteRole(t *testing.T) {
	resourceName := "luminate_site_role.site-role"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSiteRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "SiteEditor")),
			},
		},
	})
}
