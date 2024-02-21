package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func testAccResourceIntegration_minimal(name string) string {
	return fmt.Sprintf(`
	resource "luminate_aws_integration" "new-integration" {
		integration_name = "%s"
	}`, name)
}

func TestAccLuminateIntegration(t *testing.T) {
	resourceName := "luminate_aws_integration.new-integration"
	var awsAccountID, awsIntegrationName string
	if awsAccountID = os.Getenv("TEST_LUMINATE_AWS_ACCOUNT_ID"); awsAccountID == "" {
		t.Skip("skipping TestAccLuminateIntegration no  luminate aws account id provided")
	}
	if awsIntegrationName = os.Getenv("AWS_INTEGRATION_NAME"); awsIntegrationName == "" {
		t.Skip("skipping TestAccLuminateIntegration no aws integration name provided")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIntegration_minimal(awsIntegrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "luminate_aws_account_id", awsAccountID),
				),
			},
		},
	})
}
