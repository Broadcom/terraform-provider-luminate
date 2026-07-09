package framework_provider

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Broadcom/terraform-provider-luminate/test_utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCollectionDataSourceTemplate = `
resource "luminate_collection" "new_collection_<RANDOM_PLACEHOLDER>" {
	name = "tfAccCollection<RANDOM_PLACEHOLDER>"
}

data "luminate_collection" "collection_<RANDOM_PLACEHOLDER>" {
	name = luminate_collection.new_collection_<RANDOM_PLACEHOLDER>.name
}
`

func TestAccLuminateDataSourceCollection(t *testing.T) {
	randNum := test_utils.GetRandomNumber()
	resourceName := fmt.Sprintf("luminate_collection.new_collection_%s", strconv.Itoa(randNum))
	dataResourceName := fmt.Sprintf("data.luminate_collection.collection_%s", strconv.Itoa(randNum))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: strings.ReplaceAll(testAccCollectionDataSourceTemplate, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataResourceName, "id", resourceName, "id"),
				),
			},
		},
	})
}
