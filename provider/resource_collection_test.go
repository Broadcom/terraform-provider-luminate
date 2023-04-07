package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const tessAccResourceCollection = `
resource "luminate_collection" "new-collection" {
  name = "tfAccCollection"
}
`

func TestAccLuminateCollection(t *testing.T) {
	resourceName := "luminate_collection.new-collection"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tessAccResourceCollection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccCollection"),
				),
			},
		},
	})
}
