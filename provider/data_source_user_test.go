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
					resource.TestCheckResourceAttr(resourceName, "user_ids.0", "e9bb7894-a6e4-44de-a2b3-ee9e5e72485a"),
				),
			},
		},
	})
}
