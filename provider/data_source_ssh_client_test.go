package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceSshClients(name string) string {
	return fmt.Sprintf(`
	data "luminate_ssh_client" "my-ssh-client" {
	  name = "%s"
	}
`, name)
}

func TestAccLuminateDataSourceSshClients(t *testing.T) {
	resourceName := "data.luminate_ssh_client.my-ssh-client"
	var sshClientName string
	if sshClientName = os.Getenv("TEST_SSH_CLIENT_NAME"); sshClientName == "" {
		t.Skip("skipping TestAccLuminateDataSourceSshClients no ssh client name provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSshClients(sshClientName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "6ddace8e-39a3-4cd3-bba7-ad26e826df5b"),
				),
			},
		},
	})
}
