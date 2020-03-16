package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
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
				Config:    testAccResourceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_ids.0", "f879c404-6f6d-427e-b483-faa2c9d5017d"),
				),
			},
		},
	})
}
