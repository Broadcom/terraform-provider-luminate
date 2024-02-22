package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func testAccResourceIntegrationBind_minimal(id string) string {
	return fmt.Sprintf(
		`
		resource "luminate_aws_integration" "new-integration" {
			integration_name = "tfAccIntegrationBind1 "
		}
		
		resource "luminate_aws_integration_bind" "new-integration-bind" {
			integration_name = "${luminate_aws_integration.new-integration.integration_name}"
			integration_id= "${luminate_aws_integration.new-integration.integration_id}"
			aws_role_arn= "arn:aws:iam::%s:user/test"
			luminate_aws_account_id= "${luminate_aws_integration.new-integration.luminate_aws_account_id}"	
			aws_external_id= "${luminate_aws_integration.new-integration.aws_external_id}"
			regions = ["1","2"]
		}`, id)
}

func TestAccLuminateIntegrationBind(t *testing.T) {
	resourceName := "luminate_aws_integration_bind.new-integration-bind"
	var awsAccountID string
	if awsAccountID = os.Getenv("TEST_AWS_ACCOUNT_ID"); awsAccountID == "" {
		t.Skip("skipping TestAccLuminateIntegrationBind no  aws account number provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIntegrationBind_minimal(awsAccountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "luminate_aws_account_id", awsAccountID),
				),
			},
		},
	})
}
