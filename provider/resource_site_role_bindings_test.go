package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"math/rand"
	"testing"
	"time"
)

func testAccLuminateSiteRoleBindings(siteName string, rand int) string {
	return fmt.Sprintf(`resource "luminate_site" "new-site" {
		name = "%s%d"
	}
	resource "luminate_site_role" "site-admin" {
		role = "SiteEditor"
		entity_id = "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"
		identity_provider_id = "local"
		site_id = "${luminate_site.new-site.id}"
	}
`, siteName, rand)
}
func TestAccLuminateSiteRoleBindings(t *testing.T) {
	const resourceName = "luminate_site_role.site-admin"
	const siteName = "siteBindings"
	rand.Seed(time.Now().UnixNano())
	randNum := 100 + rand.Intn(100)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLuminateSiteRoleBindings(siteName, randNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role", "SiteEditor"),
					resource.TestCheckResourceAttr(resourceName, "entity_id", "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"),
				),
			},
		},
	})
}
