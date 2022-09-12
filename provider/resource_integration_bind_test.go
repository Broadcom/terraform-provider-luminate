package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testAccResourceIntegrationBind_minimal = `
resource "luminate_aws_integration" "new-integration" {
	integration_name = "tfAccIntegrationBind1 "
}

resource "luminate_aws_integration_bind" "new-integration-bind" {
	integration_name = "${luminate_aws_integration.new-integration.integration_name}"
	integration_id= "${luminate_aws_integration.new-integration.integration_id}"
	aws_role_arn= "arn:aws:iam::880793335152:user/test"
	luminate_aws_account_id= "${luminate_aws_integration.new-integration.luminate_aws_account_id}"	
	aws_external_id= "${luminate_aws_integration.new-integration.aws_external_id}"
	regions = ["1","2"]
}
`

func TestAccLuminateIntegrationBind(t *testing.T) {
	resourceName := "luminate_aws_integration_bind.new-integration-bind"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIntegrationBind_minimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "luminate_aws_account_id", "670797135152"),
				),
			},
		},
	})
}
