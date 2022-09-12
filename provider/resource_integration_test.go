package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testAccResourceIntegration_minimal = `
resource "luminate_aws_integration" "new-integration" {
	integration_name = "tfAccIntegration"
}
`

func TestAccLuminateIntegration(t *testing.T) {
	resourceName := "luminate_aws_integration.new-integration"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIntegration_minimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "luminate_aws_account_id", "670797135152"),
				),
			},
		},
	})
}
