package provider

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccSSHApplication_minimal(rand int) string {
	return fmt.Sprintf(
		`
resource "luminate_site" "new-site" {
   name = "tfAccSiteSSH%d"
}

resource "luminate_ssh_application" "new-ssh-application" {
	site_id = "${luminate_site.new-site.id}"
	name = "tfAccSSH%d"
	internal_address = "tcp://127.0.0.2"
 	icon = "iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII="
}
`, rand, rand)
}

func testAccSSHApplication_options(rand int) string {
	return fmt.Sprintf(
		`
resource "luminate_site" "new-site" {
	name = "tfAccSiteSSH%d"
}

resource "luminate_ssh_application" "new-ssh-application" {
	site_id = "${luminate_site.new-site.id}"
	name = "tfAccSSHUpd%d"
	internal_address = "tcp://127.0.0.5"
}
`, rand, rand)
}

func testAccSSHApplication_collection(rand int) string {
	return fmt.Sprintf(
		`
resource "luminate_site" "new-site-collection" {
	name = "tfAccSiteCollection%d"
}

resource "luminate_collection" "new-collection" {
	name = "tfAccCollectionForApp%d"
}
resource "luminate_collection_site_link" "new-collection-site-link" {
	site_id = "${luminate_site.new-site-collection.id}"
	collection_ids = sort(["${luminate_collection.new-collection.id}"])
}

resource "luminate_ssh_application" "new-ssh-application-collection" {
	site_id = "${luminate_site.new-site-collection.id}"
	collection_id = "${luminate_collection.new-collection.id}"
	name = "tfAccSSHWithCollection"
	internal_address = "tcp://127.0.0.5"
 	depends_on = [luminate_collection_site_link.new-collection-site-link]
}`, rand, rand)
}

func TestAccLuminateSSHApplication(t *testing.T) {
	resourceName := "luminate_ssh_application.new-ssh-application"
	resourceNameCollection := "luminate_ssh_application.new-ssh-application-collection"

	randNum := 100 + rand.Intn(100)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHApplication_minimal(randNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfAccSSH%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "visible", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "internal_address", "tcp://127.0.0.2"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tfaccssh.ssh.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tfaccssh.ssh.%s", testAccDomain)),
				),
			},
			{
				Config: testAccSSHApplication_options(randNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfAccSSHUpd%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "internal_address", "tcp://127.0.0.5"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tfaccssh.ssh.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tfaccssh.ssh.%s", testAccDomain)),
				),
			},
			{
				Config: testAccSSHApplication_collection(randNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", fmt.Sprintf("tfAccSSHWithCollection%d", randNum))),
			},
		},
	})
}
