package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccTCPApplication_minimal = `
resource "luminate_site" "new-site" {
    name = "tfAccSite"
}
resource "luminate_tcp_application" "new-tcp-application" {
  name = "tfAccTCP"
  site_id = "${luminate_site.new-site.id}"
  target {
    address = "127.0.0.1"
    ports = ["8080"]
  }
}
`

const testAccTCPApplication_with_collection = `
resource "luminate_site" "new-site" {
    name = "tfAccSite"
}
resource "luminate_collection" "new-collection" {
	name = "tfAccCollectionForApp"
}
resource "luminate_collection_site_link" "new-collection-site-link" {
	site_id = "${luminate_site.new-site.id}"
	collection_ids = sort(["${luminate_collection.new-collection.id}"])
}
resource "luminate_tcp_application" "new-tcp-application-collection" {
  name = "tfAccTCPWithCollection"
  site_id = "${luminate_site.new-site.id}"
  collection_id = "${luminate_collection.new-collection.id}"
  target {
    address = "127.0.0.1"
    ports = ["8080"]
  }
 depends_on = [luminate_collection_site_link.new-collection-site-link]
}
`

func TestAccLuminateTCPApplication(t *testing.T) {
	resourceName := "luminate_tcp_application.new-tcp-application"
	resourceNameCollection := "luminate_tcp_application.new-tcp-application-collection"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTCPApplication_with_collection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", "tfAccTCPWithCollection")),
			},
			{
				Config: testAccTCPApplication_minimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccTCP"),
					resource.TestCheckResourceAttr(resourceName, "visible", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tfacctcp.tcp.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tfacctcp.tcp.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "target.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target.0.address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "target.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target.0.ports.0", "8080"),
				),
			},
		},
	})
}
