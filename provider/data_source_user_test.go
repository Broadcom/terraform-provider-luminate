package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccResourceUser = `
	data "luminate_user"  "my-users" {
		identity_provider_id = "local"
		users = ["support.admin@babookenv.luminate-ci.com"]
	}
`

func TestAccLuminateDataSourceUser(t *testing.T) {
	resourceName := "data.luminate_user.my-users"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "users.0", "support.admin@babookenv.luminate-ci.com"),
				),
			},
		},
	})
}
