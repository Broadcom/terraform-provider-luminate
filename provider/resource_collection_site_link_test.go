package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const tessAccResourceCollectionSiteLink = `
resource "luminate_site" "new-site" {
	name = "tfAccCollectionSiteLink"
}
resource "luminate_collection_site_link" "new-collection-site-link" {
  links {
    site_id = luminate_site.new-site.id
    collection_id = "7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5"
  }
}
`

func TestAccLuminateCollectionSiteLink(t *testing.T) {
	resourceName := "luminate_collection_site_link.new-collection-site-link"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tessAccResourceCollectionSiteLink,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "links.0.collection_id", "7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5"),
				),
			},
		},
	})
}
