package provider

import (
	"github.com/hashicorp/terraform/helper/resource"

	"testing"
)

const testAccResourceSite_minimal = `
resource "luminate_site" "new-site" {
	name = "tfAccSite"
	region = "us-east4"
}
`

const testAccResourceSite_options = `
resource "luminate_site" "new-site" {
	name = "tfAccSiteOpt"
	region = "us-east4"
	mute_health_notification = "true"
	kubernetes_persistent_volume_name = "K8SVolume"
	kerberos {
		domain = "domain.com"
		kdc_address = "kdc_address"
		keytab_pair = "keytab_pair"
	}
}
`

func TestAccLuminateSite(t *testing.T) {
	resourceName := "luminate_site.new-site"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSite_minimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccSite"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east4"),
					resource.TestCheckResourceAttr(resourceName, "mute_health_notification", "false"),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_persistent_volume_name", ""),
					resource.TestCheckResourceAttr(resourceName, "kerberos.#", "0"),
				),
			},
			{
				Config: testAccResourceSite_options,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccSiteOpt"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east4"),
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
