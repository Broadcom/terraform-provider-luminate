package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourceUser = `
	data "luminate_user"  "my-users" {
		identity_provider_id = "local"
		users = ["tf-user@terraformat.luminatesec.com"]
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
					resource.TestCheckResourceAttr(resourceName, "user_ids.0", "f75f45b8-d10d-4aa6-9200-5c6d60110430"),
				),
			},
		},
	})
}
