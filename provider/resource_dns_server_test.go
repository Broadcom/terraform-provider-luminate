package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testDNSServer = `
resource "luminate_site" "new-site" {
	name = "tfAccSite"
}
resource "luminate_dns_server" "new-dns" {
	site_id = "${luminate_site.new-site.id}"
	name = "testDNS"
	internal_address = "udp://10.0.0.1:53"
	dns_settings {
		domain_suffixes = ["company.com"]
	}
	visible = false
}
`

func TestAccLuminateDNSServer(t *testing.T) {
	resourceName := "luminate_dns_server.new-dns"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:  testDNSServer,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testDNS"),
					resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "internal_address", "udp://10.0.0.1:53"),
				),
			},
		},
	})
}
