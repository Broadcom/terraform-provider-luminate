package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testAccResourceAwsIntegration = `
	data "luminate_aws_integration" "my-aws_integration" {
	  integration_name = "terraform-acceptance"
	}
`

func TestAccLuminateDataSourceAwsIntegration(t *testing.T) {
	resourceName := "data.luminate_aws_integration.my-aws_integration"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:    testAccResourceAwsIntegration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
