package serial_tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceAwsIntegration(name string) string {
	return fmt.Sprintf(`
	data "luminate_aws_integration" "my-aws_integration" {
	  integration_name = "%s"
	}
`, name)
}

func TestAccLuminateDataSourceAwsIntegration_Serial(t *testing.T) {
	resourceName := "data.luminate_aws_integration.my-aws_integration"
	var integrationName string
	if integrationName = os.Getenv("TEST_AWS_INTEGRATION_NAME"); integrationName == "" {
		t.Error("stopping TestAccLuminateDataSourceAsIntegration, no integration name provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAwsIntegration(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
