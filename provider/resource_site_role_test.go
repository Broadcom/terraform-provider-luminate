package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const testSiteRole = `
	resource "luminate_site" "site" {
		name = "siteToBeAssign"
	} 
	resource "luminate_site_role" "site-editor" {
		role_type = "SiteEditor"
		identity_provider_id =  "local"
		entity_id = "f75f45b8-d10d-4aa6-9200-5c6d60110430"
		entity_type = "User"
		site_id = "${luminate_site.site.id}"
	}
	resource "luminate_site_role" "site-deployer" {
		role_type = "SiteConnectorDeployer"
		identity_provider_id =  "local"
		entity_id = "f75f45b8-d10d-4aa6-9200-5c6d60110430"
		entity_type = "User"
		site_id = "${luminate_site.site.id}"
	}
`

func TestAccLuminateSiteRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSiteRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("luminate_site_role.site-editor", "role_type", "SiteEditor"),
					resource.TestCheckResourceAttr("luminate_site_role.site-deployer", "role_type", "SiteConnectorDeployer"),
				),
			},
		},
	})
}
