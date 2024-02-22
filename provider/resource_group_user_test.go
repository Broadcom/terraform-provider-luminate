package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"

	"testing"
)

func testAccResourceGroupUsers(groupName string, username string) string {
	return fmt.Sprintf(`
			data "luminate_group" "my-groups" {
				identity_provider_id = "local"
				groups = ["%s"]
			}
			
			data "luminate_user" "my-users" {
				identity_provider_id = "local"
				users = ["%s"]
			}
			
			resource "luminate_group_user" "new_group_membership" {
				group_id = "${data.luminate_group.my-groups.group_ids.0}"
				user_id = "${data.luminate_user.my-users.user_ids.0}"
			}`, groupName, username)
}

func TestAccLuminateGroupUser(t *testing.T) {
	resourceName := "luminate_group_user.new_group_membership"
	var username, groupName string
	if username = os.Getenv("TEST_USERNAME"); username == "" {
		t.Skip("skipping TestAccLuminateDataSourceUser no username provided")
	}
	if groupName = os.Getenv("TEST_GROUP_NAME"); groupName == "" {
		t.Skip("skipping TestAccLuminateDataSourceUser no  groupName provided")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupUsers(groupName, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
