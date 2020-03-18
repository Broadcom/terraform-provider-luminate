package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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

func TestAccLuminateTCPApplication(t *testing.T) {
	resourceName := "luminate_tcp_application.new-tcp-application"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTCPApplication_minimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccTCP"),
					resource.TestCheckResourceAttr(resourceName, "visible", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tfacctcp.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tfacctcp.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "target.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target.0.address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "target.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target.0.ports.0", "8080"),
				),
			},
		},
	})
}
