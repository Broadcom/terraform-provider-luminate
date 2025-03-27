package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

const testAccResourceCollection = `
resource "luminate_collection" "new-collection" {
  name = "tfAccCollection<RANDOM_PLACEHOLDER>"
}
`

const testAccResourceCollectionUpdate = `
resource "luminate_collection" "new-collection" {
  name = "tfAccCollectionUpdate<RANDOM_PLACEHOLDER>"
}
`

func TestAccLuminateCollection(t *testing.T) {
	resourceName := "luminate_collection.new-collection"
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: strings.ReplaceAll(testAccResourceCollection, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfAccCollection%d", randNum)),
				),
			},
			{
				Config: strings.ReplaceAll(testAccResourceCollectionUpdate, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfAccCollectionUpdate%d", randNum)),
				),
			},
		},
	})
}
