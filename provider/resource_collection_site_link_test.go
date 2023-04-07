package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"math/rand"
	"testing"
	"time"
)

const tessAccResourceCollectionSiteLinkUpdate = `
resource "luminate_site" "new-site" {
	name = "tfAccCollectionSiteLink"
}
resource "luminate_collection" "new-collection" {
	name = "tfAccCollection"
}
resource "luminate_collection_site_link" "new-collection-site-link" {
	site_id = "${luminate_site.new-site.id}"
	collection_ids = ["${luminate_collection.new-collection.id}"]
}
`

func TestAccLuminateCollectionSiteLink(t *testing.T) {
	resourceName := "luminate_collection_site_link.new-collection-site-link"
	rand.Seed(time.Now().UnixNano())
	randNum := 100 + rand.Intn(100)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:  testAccResourceCollectionSiteCreate("tfAccCollection", randNum),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "collection_ids.0", "7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5"),
					resource.TestCheckResourceAttr(resourceName, "collection_ids.#", "1"),
				),
			},
			{
				Config:  testAccResourceCollectionSiteUpdateSwitch("tfAccCollection", randNum),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "collection_ids.#", "1"),
				),
			},
			{
				Config: testAccResourceCollectionSiteUpdateOnlyAdd("tfAccCollection", randNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "collection_ids.#", "2"),
				),
			},
			//{
			//	Config: testAccTemp(randNum),
			//},
		},
	})
}

func testAccResourceCollectionSiteCreate(collectionName string, rand int) string {
	return fmt.Sprintf(`
				resource "luminate_site" "new-site" {
					name = "tfAccCollectionSiteLink%d"
				}
				resource "luminate_collection" "new-collection" {
					name = "%s%d"
				}
				resource "luminate_collection_site_link" "new-collection-site-link" {
					site_id = "${luminate_site.new-site.id}"
					collection_ids = ["7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5"]
				}
			`, rand, collectionName, rand)
}

func testAccResourceCollectionSiteUpdateSwitch(collectionName string, rand int) string {
	return fmt.Sprintf(`
				resource "luminate_site" "new-site" {
					name = "tfAccCollectionSiteLink%d"
				}
				resource "luminate_collection" "new-collection" {
					name = "%s%d"
				}
				resource "luminate_collection_site_link" "new-collection-site-link" {
					site_id = "${luminate_site.new-site.id}"
					collection_ids = ["${luminate_collection.new-collection.id}"]
				}
			`, rand, collectionName, rand)
}

func testAccResourceCollectionSiteUpdateOnlyAdd(collectionName string, rand int) string {
	return fmt.Sprintf(`
				resource "luminate_site" "new-site" {
					name = "tfAccCollectionSiteLink%d"
				}
				resource "luminate_collection" "new-collection" {
					name = "%s%d"
				}
				resource "luminate_collection_site_link" "new-collection-site-link" {
					site_id = "${luminate_site.new-site.id}"
					collection_ids = ["7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5", "${luminate_collection.new-collection.id}"]
				}
			`, rand, collectionName, rand)
}

func testAccTemp(rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
		name = "tfSiteCollection%d"
	}
	
	resource "luminate_collection" "new-collection" {
	  name = "tfAccCollection%d"
	}
	
	resource "luminate_collection_site_link" "new-collection-site-link" {
		site_id = "${luminate_site.new-site.id}"
		collection_ids = [ "7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5"]
	}
	
	
	resource "luminate_web_application" "new-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccApplication%d"
		internal_address = "http://127.0.0.1:8080"
	}
	
	resource "luminate_web_access_policy" "new-web-access-policy" {
		enabled = "true"
		name =  "resourceWebAccessPolicy%d"
		identity_provider_id = "local"
	
		user_ids = ["24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"]
		applications = ["${luminate_web_application.new-application.id}"]
	}
`, rand, rand, rand, rand)
}
