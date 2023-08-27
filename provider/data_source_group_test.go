package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourceGroup = `
	data "luminate_group"  "my-groups" {
		identity_provider_id = "local"
		groups = ["tf-acceptance"]
	}
`

func TestAccLuminateDataSourceGroup(t *testing.T) {
	resourceName := "data.luminate_group.my-groups"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "groups.0", "tf-acceptance"),
				),
			},
		},
	})
}
