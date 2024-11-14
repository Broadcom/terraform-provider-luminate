package provider

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDNSServerResiliency = fmt.Sprintf(`
resource "luminate_site" "new-site" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
}

resource "luminate_collection_site_link" "new-collection-site-link" {
	site_id = "${luminate_site.new-site.id}"
	collection_ids = ["7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5"]
}

resource "luminate_dns_group_resiliency" "new-dns-group" {
	name = "testDNSGroupResiliency<RANDOM_PLACEHOLDER>"
	domain_suffixes = ["%s"]
	send_notifications = true
}

resource "luminate_dns_server_resiliency" "new-dns-server-resiliency" {
	name = "testDNSServerResiliency<RANDOM_PLACEHOLDER>"
	site_id = "${luminate_site.new-site.id}"
	group_id = "${luminate_dns_group_resiliency.new-dns-group.id}"
	internal_address = "udp://<RANDOM_PLACEHOLDER>.<RANDOM_PLACEHOLDER>.<RANDOM_PLACEHOLDER>.<RANDOM_PLACEHOLDER>:63"
}

`, "otherdomains<RANDOM_PLACEHOLDER>.com")

func TestAccLuminateDNSServerResiliency(t *testing.T) {
	resourceName := "luminate_dns_server_resiliency.new-dns-server-resiliency"
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: strings.ReplaceAll(testDNSServerResiliency, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testDNSServerResiliency%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "internal_address", strings.ReplaceAll("udp://<RANDOM_PLACEHOLDER>.<RANDOM_PLACEHOLDER>.<RANDOM_PLACEHOLDER>.<RANDOM_PLACEHOLDER>:63", "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum))),
				),
			},
		},
	})
}
