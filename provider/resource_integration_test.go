package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

const testAccResourceIntegration_minimal = `
resource "luminate_aws_integration" "new-integration" {
	integration_name = "tfAccIntegration"
}
`

func TestAccLuminateIntegration(t *testing.T) {
	resourceName := "luminate_aws_integration.new-integration"
	var awsAccountID string
	if awsAccountID = os.Getenv("TEST_AWS_ACCOUNT_ID"); awsAccountID == "" {
		t.Skip("skipping TestAccLuminateIntegration no luminate aws account id provided")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIntegration_minimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "luminate_aws_account_id", awsAccountID),
				),
			},
		},
	})
}
