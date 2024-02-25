package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func testSiteRole(name string) string {
	return fmt.Sprintf(`
	resource "luminate_site" "site" {
		name = "siteToBeAssign"
	} 
	resource "luminate_site_role" "site-editor" {
		role_type = "SiteEditor"
		identity_provider_id =  "local"
		entity_id = "%s"
		entity_type = "User"
		site_id = "${luminate_site.site.id}"
	}
	resource "luminate_site_role" "site-deployer" {
		role_type = "SiteConnectorDeployer"
		identity_provider_id =  "local"
		entity_id = "%s"
		entity_type = "User"
		site_id = "${luminate_site.site.id}"
	}
`, name, name)
}

func TestAccLuminateSiteRole(t *testing.T) {
	var userID string
	if userID = os.Getenv("TEST_USER_ID"); userID == "" {
		t.Error("skipping TestAccLuminateSiteRole no user id provided")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSiteRole(userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("luminate_site_role.site-editor", "role_type", "SiteEditor"),
					resource.TestCheckResourceAttr("luminate_site_role.site-deployer", "role_type", "SiteConnectorDeployer"),
				),
			},
		},
	})
}
