package provider

import (
	"github.com/hashicorp/terraform/helper/resource"

	"testing"
)

const testAccResourceGroupUsers = `
data "luminate_group" "my-groups" {
	identity_provider_id = "local"
	groups = ["tf-acceptance"]
}

data "luminate_user" "my-users" {
	identity_provider_id = "local"
	users = ["support.admin@babookenv.luminate-ci.com"]
}

resource "luminate_group_user" "new_group_membership" {
	group_id = "${data.luminate_group.my-groups.group_ids.0}"
	user_id = "${data.luminate_user.my-users.user_ids.0}"
}
`

func TestAccLuminateGroupUser(t *testing.T) {
	resourceName := "luminate_group_user.new_group_membership"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupUsers,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
