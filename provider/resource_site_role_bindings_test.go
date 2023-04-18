package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testAccLuminateSiteRoleBindings = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite"
	}
	resource "luminate_site_role" "site-admin" {
		role = "SiteEditor"
		entity_id = "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"
		identity_provider_id = "local"
		site_id = "${luminate_site.new-site.id}"
	}
`

func TestAccLuminateRoleBindings(t *testing.T) {
	const resourceName = "luminate_site_role.site-admin"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:  testAccLuminateSiteRoleBindings,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role", "SiteEditor"),
					resource.TestCheckResourceAttr(resourceName, "entity_id", "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"),
				),
			},
		},
	})
}
