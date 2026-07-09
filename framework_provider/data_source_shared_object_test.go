package framework_provider

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Broadcom/terraform-provider-luminate/test_utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccSharedObjectDataSourceTemplate = `
data "luminate_shared_object" "shared_object_<RANDOM_PLACEHOLDER>" {
	name = "tfAccSharedObject<RANDOM_PLACEHOLDER>"
	type = "ISOLATION_PROFILE"
}
`

func TestAccLuminateDataSourceSharedObject(t *testing.T) {
	randNum := test_utils.GetRandomNumber()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config:      strings.ReplaceAll(testAccSharedObjectDataSourceTemplate, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				ExpectError: regexp.MustCompile(".*Shared object not found"),
			},
		},
	})
}
