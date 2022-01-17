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
					resource.TestCheckResourceAttr(resourceName, "id", "6ddace8e-39a3-4cd3-bba7-ad26e826df5b"),
				),
			},
		},
	})
}
