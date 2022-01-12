package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccSSHApplication_minimal = `
resource "luminate_site" "new-site" {
   name = "tfAccSite"
}

resource "luminate_ssh_application" "new-ssh-application" {
	site_id = "${luminate_site.new-site.id}"
	name = "tfAccSSH"
	internal_address = "tcp://127.0.0.2"
}
`

const testAccSSHApplication_options = `
resource "luminate_site" "new-site" {
	name = "tfAccSite"
}

resource "luminate_ssh_application" "new-ssh-application" {
	site_id = "${luminate_site.new-site.id}"
	name = "tfAccSSHUpd"
	internal_address = "tcp://127.0.0.5"
}
`

func TestAccLuminateSSHApplication(t *testing.T) {
	resourceName := "luminate_ssh_application.new-ssh-application"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHApplication_minimal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccSSH"),
					resource.TestCheckResourceAttr(resourceName, "visible", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "internal_address", "tcp://127.0.0.2"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tfaccssh.ssh.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tfaccssh.ssh.%s", testAccDomain)),
				),
			},
			{
				Config: testAccSSHApplication_options,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tfAccSSHUpd"),
					resource.TestCheckResourceAttr(resourceName, "internal_address", "tcp://127.0.0.5"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tfaccssh.ssh.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tfaccssh.ssh.%s", testAccDomain)),
				),
			},
		},
	})
}
