package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourceIdentityProvider = `
	data "luminate_identity_provider" "my-identity-provider" {
	  identity_provider_name = "local"
	}
`

func TestAccLuminateDataSourceIdentityProvider(t *testing.T) {
	resourceName := "data.luminate_identity_provider.my-identity-provider"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIdentityProvider,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identity_provider_id", "local"),
				),
			},
		},
	})
}
