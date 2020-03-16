package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testAccResourceUser = `
	data "luminate_user"  "my-users" {
		identity_provider_id = "local"
		users = ["tf-user@tfacc.luminatesec.com"]
	}
`

func TestAccLuminateDataSourceUser(t *testing.T) {
	resourceName := "data.luminate_user.my-users"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:    testAccResourceUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_ids.0", "f04d9234-3482-48b0-b56b-d562a5d90f26"),
				),
			},
		},
	})
}
