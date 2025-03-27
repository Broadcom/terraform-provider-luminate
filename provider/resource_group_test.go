package provider

import (
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const GroupNameRandomPartSize = 7

func testGroupCreate(idpID string) string {
	randomGroupName := "testGroup" + utils.GenerateRandomString(GroupNameRandomPartSize)
	return fmt.Sprintf(`
    resource "luminate_resources_group" "new-group" {
	name = "%s"
	identity_provider_id = "%s"
		}
	`, randomGroupName, idpID)
}

func TestGroupCreate(t *testing.T) {
	resourceName := "luminate_resources_group.new-group"
	var idpID string
	if idpID = os.Getenv("TEST_IDP_ID"); idpID == "" {
		t.Error("stopping TestGroupCreate no idpID provided")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: testGroupCreate(idpID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", createRegExpForNamePrefix("testGroup")),
				),
			},
		},
	})
}
