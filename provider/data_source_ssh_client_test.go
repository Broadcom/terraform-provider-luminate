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
	var sshClientName, sshClientID string
	if sshClientName = os.Getenv("TEST_SSH_CLIENT_NAME"); sshClientName == "" {
		t.Error("stopping TestAccLuminateDataSourceSshClients no ssh client name provided")
	}
	if sshClientID = os.Getenv("TEST_SSH_CLIENT_ID"); sshClientID == "" {
		t.Error("stopping TestAccLuminateDataSourceSshClients no ssh client id provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSshClients(sshClientName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", sshClientID),
				),
			},
		},
	})
}
