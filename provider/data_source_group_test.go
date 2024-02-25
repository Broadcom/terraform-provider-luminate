package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func getTestAccResourceGroup(name string) string {
	return fmt.Sprintf(`
	data "luminate_group"  "my-groups" {
		identity_provider_id = "local"
		groups = ["%s"]
	}
`, name)

}

func TestAccLuminateDataSourceGroup(t *testing.T) {
	resourceName := "data.luminate_group.my-groups"
	var groupName string
	if groupName = os.Getenv("TEST_GROUP_NAME"); groupName == "" {
		t.Error("skipping TestAccLuminateDataSourceGroup no group name provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceGroup(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "groups.0", groupName),
				),
			},
		},
	})
}
