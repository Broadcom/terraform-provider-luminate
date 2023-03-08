package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccResourceSshClients = `
	data "luminate_ssh_client" "my-ssh-client" {
	  name = "tf-at-ssh-client"
	}
`

func TestAccLuminateDataSourceSshClients(t *testing.T) {
	resourceName := "data.luminate_ssh_client.my-ssh-client"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSshClients,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "324bad9c-c6ea-4c3c-874f-cdbd3d734555"),
				),
			},
		},
	})
}
