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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_ids.0", "3b61849d-f08d-42d3-a158-da1a53cd2ac6"),
				),
			},
		},
	})
}
