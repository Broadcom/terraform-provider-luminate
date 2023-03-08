package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccWebApplication_minimal = `
resource "luminate_site" "new-site" {
   name = "tfAccSite"
}
resource "luminate_web_application" "new-application" {
 site_id = "${luminate_site.new-site.id}"
 name = "tfAccApplication"
 internal_address = "http://127.0.0.1:8080"
}
`

const testAccWebApplication_options = `
resource "luminate_site" "new-site" {
   name = "tfAccSite"
}
resource "luminate_web_application" "new-application" {
 site_id = "${luminate_site.new-site.id}"
 name = "tfAccApplicationUpd"
 internal_address = "http://127.0.0.1:80"
 custom_root_path = "/testAcc"
}
`

func TestAccLuminateApplication(t *testing.T) {
	resourceName := "luminate_web_application.new-application"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWebApplication_minimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccApplication"),
					resource.TestCheckResourceAttr(resourceName, "visible", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "internal_address", "http://127.0.0.1:8080"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("https://tfaccapplication.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("https://tfaccapplication.%s", testAccDomain)),
				),
			},
			{
				Config: testAccWebApplication_options,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccApplicationUpd"),
					resource.TestCheckResourceAttr(resourceName, "internal_address", "http://127.0.0.1:80"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("https://tfaccapplication.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("https://tfaccapplication.%s", testAccDomain)),
				),
			},
		},
	})
}
