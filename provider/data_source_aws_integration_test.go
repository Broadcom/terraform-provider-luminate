package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourceAwsIntegration = `
	data "luminate_aws_integration" "my-aws_integration" {
	  integration_name = "terraform-test"
	}
`

func TestAccLuminateDataSourceAwsIntegration(t *testing.T) {
	resourceName := "data.luminate_aws_integration.my-aws_integration"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAwsIntegration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
