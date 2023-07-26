package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourceConnector = `
resource "luminate_site" "new-site" {
name = "tfAccConnector"
}

resource "luminate_connector" "new-connector" {
site_id = "${luminate_site.new-site.id}"
name = "connector"
type = "linux"
}
`

func TestAccLuminateConnector(t *testing.T) {
	resourceName := "luminate_connector.new-connector"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnector,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "connector"),
					resource.TestCheckResourceAttr(resourceName, "type", "linux"),
				),
			},
		},
	})
}
