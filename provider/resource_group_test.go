package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testGroupCreate(idpID string) string {
	return fmt.Sprintf(`
    resource "luminate_resources_group" "new-group" {
	name = "testGroup"
	identity_provider_id = "%s"
		}
	`, idpID)
}

func TestGroupCreate(t *testing.T) {
	resourceName := "luminate_resources_group.new-group"
	var idpID string
	if idpID = os.Getenv("TEST_IDP_ID"); idpID == "" {
		t.Skip("won't fail since not merged to AT yet")
		// todo when AT is ready enable error & remove skip
		//t.Error("stopping TestGroupCreate no idpID provided")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGroupCreate(idpID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testGroup"),
				),
			},
		},
	})
}
