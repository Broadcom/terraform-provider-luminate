package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestAccLuminateCollectionSiteLink(t *testing.T) {
	resourceNameSiteLinks := "luminate_collection_site_link.new-collection-site-link"
	resourceNameSite := "luminate_site.new-site"
	resourceNameCollection := "luminate_collection.new-collection"
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
					resource.TestCheckResourceAttr(resourceNameSiteLinks, "collection_ids.0", "7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5"),
					resource.TestCheckResourceAttr(resourceNameSiteLinks, "collection_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameSiteLinks, "site_id", resourceNameSite, "id"),
				),
			},
			{
				Config:  testAccResourceCollectionSiteUpdateSwitch("tfAccCollection", randNum),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSiteLinks, "collection_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameSiteLinks, "collection_ids.0", resourceNameCollection, "id"),
					resource.TestCheckResourceAttrPair(resourceNameSiteLinks, "site_id", resourceNameSite, "id"),
				),
			},
			{
				Config: testAccResourceCollectionSiteUpdateOnlyAddOne("tfAccCollection", randNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSiteLinks, "collection_ids.#", "2"),
					resource.TestCheckResourceAttrPair(resourceNameSiteLinks, "site_id", resourceNameSite, "id"),
				),
			},
			{
				Config: testAccResourceCollectionSiteUpdateSwitchOrder("tfAccCollection", randNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSiteLinks, "collection_ids.#", "2"),
					resource.TestCheckResourceAttrPair(resourceNameSiteLinks, "site_id", resourceNameSite, "id"),
				),
			},
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

func testAccResourceCollectionSiteUpdateOnlyAddOne(collectionName string, rand int) string {
	return fmt.Sprintf(`
				resource "luminate_site" "new-site" {
					name = "tfAccCollectionSiteLink%d"
				}
				resource "luminate_collection" "new-collection" {
					name = "%s%d"
				}
				resource "luminate_collection_site_link" "new-collection-site-link" {
					site_id = "${luminate_site.new-site.id}"
					collection_ids = sort(["7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5", "${luminate_collection.new-collection.id}"])
				}
			`, rand, collectionName, rand)
}

func testAccResourceCollectionSiteUpdateSwitchOrder(collectionName string, rand int) string {
	return fmt.Sprintf(`
				resource "luminate_site" "new-site" {
					name = "tfAccCollectionSiteLink%d"
				}
				resource "luminate_collection" "new-collection" {
					name = "%s%d"
				}
				resource "luminate_collection_site_link" "new-collection-site-link" {
					site_id = "${luminate_site.new-site.id}"
					collection_ids = sort(["${luminate_collection.new-collection.id}", "7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5"])
				}
			`, rand, collectionName, rand)
}

func TestUniqueValues_Ok(t *testing.T) {
	leftSlice := []string{"a", "b", "c"}
	rightSlice := []string{"d", "e", "f"}
	left, right := GetUniqueValues(leftSlice, rightSlice)
	assert.Len(t, left, 3)
	assert.Len(t, right, 3)

	leftSlice = []string{"a", "b", "c"}
	rightSlice = []string{"a", "b", "c"}

	left, right = GetUniqueValues(leftSlice, rightSlice)

	assert.Len(t, left, 0)
	assert.Len(t, right, 0)

	leftSlice = []string{"a", "b", "c"}
	rightSlice = []string{"a", "b", "c", "d"}

	left, right = GetUniqueValues(leftSlice, rightSlice)

	assert.Len(t, left, 1)
	assert.Len(t, right, 0)
}
