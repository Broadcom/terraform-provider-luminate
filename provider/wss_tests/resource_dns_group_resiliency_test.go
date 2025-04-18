package wss_tests

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDNSGroupResiliency = fmt.Sprintf(`resource "luminate_dns_group_resiliency" "new-dns-group" {
	name = "testDNSGroupResiliency<RANDOM_PLACEHOLDER>"
	domain_suffixes = ["%s"]
	send_notifications = true
	is_enabled = true
}
`, "somedomain<RANDOM_PLACEHOLDER>.com")

func TestAccLuminateDNSGroupResiliency(t *testing.T) {
	runWSSTests := os.Getenv("RUN_WSS_TESTS")
	if runWSSTests != "true" {
		t.Skip("no wss tests needed")
	}
	resourceName := "luminate_dns_group_resiliency.new-dns-group"
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: strings.ReplaceAll(testDNSGroupResiliency, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testDNSGroupResiliency%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "send_notifications", "true"),
					resource.TestCheckResourceAttr(resourceName, "domain_suffixes.0", fmt.Sprintf("somedomain%s.com", strconv.Itoa(randNum))),
				),
			},
		},
	})
}
