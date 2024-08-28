package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testUserCreate_minimal = `
data "luminate_user" "bala"{
			identity_provider_id = "9748fcbe-7eac-4dc9-9809-127ca4f889ba"
			users = ["amir1@gmail.com","amir2@gmail.com"]
}

resource "luminate_resources_saml_delete_user" "new-users" {
	 identity_provider_id= "${data.luminate_user.bala.identity_provider_id}"
	 users= "${data.luminate_user.bala.users}"
}
`

func TestUserDelete(t *testing.T) {
	resourceName := "luminate_resources_saml_delete_user.new-users"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserCreate_minimal,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testUser"),
					//resource.TestCheckResourceAttr(resourceName, "visible", "true"),
					//resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					//resource.TestCheckResourceAttr(resourceName, "internal_address", "tcp://127.0.0.2:3389"),
					//resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tfaccrdp.rdp.%s", testAccDomain)),
					//resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tfaccrdp.rdp.%s", testAccDomain)),
				),
			},
		},
	})
}

//func testAccResourceGroupUsers(groupName string, username string) string {
//	return fmt.Sprintf(`
//			data "luminate_group" "my-groups" {
//				identity_provider_id = "local"
//				groups = ["%s"]
//			}
//
//			data "luminate_user" "my-users" {
//				identity_provider_id = "local"
//				users = ["%s"]
//			}
//
//			resource "luminate_group_user" "new_group_membership" {
//				group_id = "${data.luminate_group.my-groups.group_ids.0}"
//				user_id = "${data.luminate_user.my-users.user_ids.0}"
//			}`, groupName, username)
//}
//
//func TestAccLuminateGroupUser(t *testing.T) {
//	resourceName := "luminate_group_user.new_group_membership"
//	var username, groupName string
//	if username = os.Getenv("TEST_USERNAME"); username == "" {
//		t.Error("stopping TestAccLuminateDataSourceUser no username provided")
//	}
//	if groupName = os.Getenv("TEST_GROUP_NAME"); groupName == "" {
//		t.Error("stopping TestAccLuminateDataSourceUser no  groupName provided")
//	}
//	resource.Test(t, resource.TestCase{
//		PreCheck:          func() { testAccPreCheck(t) },
//		ProviderFactories: newTestAccProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccResourceGroupUsers(groupName, username),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//				),
//			},
//		},
//	})
//}
