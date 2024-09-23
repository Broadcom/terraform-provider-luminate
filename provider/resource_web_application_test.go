package provider

import (
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccWebApplication_minimal(rand int) string {
	return fmt.Sprintf(
		`
resource "luminate_site" "new-site" {
   name = "tfAccSiteForWebApp%d"
}
resource "luminate_web_application" "new-application" {
 site_id = "${luminate_site.new-site.id}"
 name = "tfAccApplication"
 internal_address = "http://127.0.0.1:8080"
 icon = "iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII="
}
`, rand)
}

func testAccWebApplication_options(rand int) string {
	return fmt.Sprintf(
		`
resource "luminate_site" "new-site" {
   name = "tfAccSiteForWebApp%d"
}
resource "luminate_web_application" "new-application" {
 site_id = "${luminate_site.new-site.id}"
 name = "tfAccApplicationUpd"
 internal_address = "http://127.0.0.1:80"
	custom_root_path = "/testAcc"
}
`, rand)
}

func testAccWebApplication_with_collection(rand int) string {
	return fmt.Sprintf(
		`
	resource "luminate_site" "new-site" {
		name = "tfAccSiteForWebApp%d"
	}
	resource "luminate_collection" "new-collection" {
		name = "tfAccCollectionForApp%d"
	}
	resource "luminate_collection_site_link" "new-collection-site-link" {
		site_id = "${luminate_site.new-site.id}"
		collection_ids = sort(["${luminate_collection.new-collection.id}"])
	}
	resource "luminate_web_application" "new-collection-application" {
		site_id = "${luminate_site.new-site.id}"
		collection_id = "${luminate_collection.new-collection.id}"
		name = "tfAccApplicationWithCollection"
		internal_address = "http://127.0.0.1:80"
		custom_root_path = "/testAcc"

 		depends_on = [luminate_collection_site_link.new-collection-site-link]
	}
	`, rand, rand)
}

func TestAccLuminateApplication(t *testing.T) {
	resourceCollectionTest := "luminate_web_application.new-collection-application"
	resourceTest := "luminate_web_application.new-application"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWebApplication_with_collection(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceCollectionTest, "name", "tfAccApplicationWithCollection"),
				),
			},
			{
				Config: testAccWebApplication_minimal(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceTest, "name", "tfAccApplication"),
					resource.TestCheckResourceAttr(resourceTest, "visible", "true"),
					resource.TestCheckResourceAttr(resourceTest, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceTest, "internal_address", "http://127.0.0.1:8080"),
					resource.TestCheckResourceAttr(resourceTest, "external_address", fmt.Sprintf("https://tfaccapplication.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceTest, "luminate_address", fmt.Sprintf("https://tfaccapplication.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceTest, "collection_id", utils.DefaultCollection),
				),
			},
			{
				Config: testAccWebApplication_options(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceTest, "name", "tfAccApplicationUpd"),
					resource.TestCheckResourceAttr(resourceTest, "internal_address", "http://127.0.0.1:80"),
					resource.TestCheckResourceAttr(resourceTest, "collection_id", utils.DefaultCollection),
					resource.TestCheckResourceAttr(resourceTest, "external_address", fmt.Sprintf("https://tfaccapplication.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceTest, "luminate_address", fmt.Sprintf("https://tfaccapplication.%s", testAccDomain)),
				),
			},
		},
	})
}
