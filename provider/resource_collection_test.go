package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testAccResourceCollection = `
resource "luminate_collection" "new-collection" {
  name = "tfAccCollection"
}
`

const testAccResourceCollectionUpdate = `
resource "luminate_collection" "new-collection" {
  name = "tfAccCollectionUpdate"
}
`

func TestAccLuminateCollection(t *testing.T) {
	resourceName := "luminate_collection.new-collection"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCollection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccCollection"),
				),
			},
			{
				Config: testAccResourceCollectionUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccCollectionUpdate"),
				),
			},
		},
	})
}
