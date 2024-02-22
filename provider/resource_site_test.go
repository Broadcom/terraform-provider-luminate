package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"

	"testing"
)

const testAccResourceSite_minimal = `
resource "luminate_site" "new-site" {
	name = "tfAccSite"
}
`

func testAccResourceSite_options(region string) string {
	return fmt.Sprintf(`resource "luminate_site" "new-site" {
	name = "tfAccSiteOpt"
	region = "%s"
	mute_health_notification = "true"
	kubernetes_persistent_volume_name = "K8SVolume"
	kerberos {
		domain = "domain.com"
		kdc_address = "kdc_address"
		keytab_pair = "keytab_pair"
	}
}`, region)
}

func TestAccLuminateSite(t *testing.T) {
	resourceName := "luminate_site.new-site"
	var region string
	if region = os.Getenv("TEST_SITE_REGION"); region == "" {
		t.Skip("skipping TestAccLuminateSite no  site provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSite_minimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccSite"),
					resource.TestCheckResourceAttr(resourceName, "mute_health_notification", "false"),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_persistent_volume_name", ""),
					resource.TestCheckResourceAttr(resourceName, "kerberos.#", "0"),
				),
			},
			{
				Config: testAccResourceSite_options(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccSiteOpt"),
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
