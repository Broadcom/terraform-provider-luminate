package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testGroupCreate_minimal = `
resource "luminate_resources_group" "new-group" {
	name = "testGroup7"
}
`

func TestGroupCreate(t *testing.T) {
	resourceName := "luminate_resources_group.new-group"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:  testGroupCreate_minimal,
				Destroy: false,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testGroup7"),
					resource.TestCheckResourceAttr(resourceName, "name", "testGroup7"),
				),
			},
		},
	})
}
