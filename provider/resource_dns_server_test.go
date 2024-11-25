package provider

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testDNSServer = `
resource "luminate_site" "new-site" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
}
resource "luminate_dns_server" "new-dns" {
	site_id = "${luminate_site.new-site.id}"
	name = "testDNS<RANDOM_PLACEHOLDER>"
	internal_address = "udp://10.0.0.1:53"
	dns_settings {
		domain_suffixes = ["company.com"]
	}
	visible = false
}
`

func TestAccLuminateDNSServer(t *testing.T) {
	resourceName := "luminate_dns_server.new-dns"
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: strings.ReplaceAll(testDNSServer, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testDNS%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "internal_address", "udp://10.0.0.1:53"),
				),
			},
		},
	})
}
