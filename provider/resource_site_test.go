package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
)

const testAccResourceSite_minimal = `
resource "luminate_site" "new-site" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
}
`

func testAccResourceSite_options(region string, randNum int) string {
	return fmt.Sprintf(`resource "luminate_site" "new-site" {
	name = "tfAccSiteOpt%d"
	region = "%s"
	mute_health_notification = "true"
	kubernetes_persistent_volume_name = "K8SVolume"
	kerberos {
		domain = "domain.com"
		kdc_address = "kdc_address"
		keytab_pair = "keytab_pair"
	}
}`, randNum, region)
}

func TestAccLuminateSite(t *testing.T) {
	resourceName := "luminate_site.new-site"
	var region string
	if region = os.Getenv("TEST_SITE_REGION"); region == "" {
		t.Error("stopping TestAccLuminateSite no  site provided")
	}
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: strings.ReplaceAll(testAccResourceSite_minimal, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfAccSite%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "mute_health_notification", "false"),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_persistent_volume_name", ""),
					resource.TestCheckResourceAttr(resourceName, "kerberos.#", "0"),
				),
			},
			{
				Config: testAccResourceSite_options(region, randNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfAccSiteOpt%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "mute_health_notification", "true"),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_persistent_volume_name", "K8SVolume"),
					resource.TestCheckResourceAttr(resourceName, "kerberos.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kerberos.0.domain", "domain.com"),
					resource.TestCheckResourceAttr(resourceName, "kerberos.0.kdc_address", "kdc_address"),
					resource.TestCheckResourceAttr(resourceName, "kerberos.0.keytab_pair", "keytab_pair"),
				),
			},
		},
	})
}
