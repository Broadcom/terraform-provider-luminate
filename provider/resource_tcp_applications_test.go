package provider

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccTCPApplication_minimal(rand int) string {
	return fmt.Sprintf(
		`
resource "luminate_site" "new-site" {
    name = "tfAccSite%d"
}
resource "luminate_tcp_application" "new-tcp-application" {
  name = "tfAccTCP%d"
  site_id = "${luminate_site.new-site.id}"
  icon = "iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII="
  target {
    address = "127.0.0.1"
    ports = ["8080"]
	port_mapping = [80]	
  }
}`, rand, rand)
}

func testAccTCPApplication_with_collection(rand int) string {
	return fmt.Sprintf(`
resource "luminate_site" "new-site" {
    name = "tfAccSite%d"
}
resource "luminate_collection" "new-collection" {
	name = "tfAccCollectionForApp%d"
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
    ports = [8080]
	port_mapping = [80]	
  }
 depends_on = [luminate_collection_site_link.new-collection-site-link]
} `, rand, rand)
}

func TestAccLuminateTCPApplication(t *testing.T) {
	resourceName := "luminate_tcp_application.new-tcp-application"
	resourceNameCollection := "luminate_tcp_application.new-tcp-application-collection"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTCPApplication_with_collection(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", "tfAccTCPWithCollection")),
			},
			{
				Config: testAccTCPApplication_minimal(100 + rand.Intn(100)),
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
