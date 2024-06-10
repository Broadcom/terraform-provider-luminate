package provider

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccRDPApplication_minimal = `
resource "luminate_site" "new-site" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
}
resource "luminate_rdp_application" "new-rdp-application" {
	site_id = "${luminate_site.new-site.id}"
	name = "tfAccRDP<RANDOM_PLACEHOLDER>"
	internal_address = "127.0.0.2"
 	icon = "iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII="
}
`
const testAccRDPApplication_changeInternalAddress = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.2:33"
	}
`

const testAccRDPApplication_changeInternalAddress_2 = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.2"
	}
`

const testAccRDPApplication_changeInternalAddress_3 = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.3:3389"
	}
`

const testAccRDPApplication_options = `
resource "luminate_site" "new-site" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
}

resource "luminate_rdp_application" "new-rdp-application" {
	site_id = "${luminate_site.new-site.id}"
	name = "tfAccRDPUpd<RANDOM_PLACEHOLDER>"
	internal_address = "tcp://127.0.0.5:126"
}
`
const testAccRDPApplication_collection = `
resource "luminate_site" "new-site-collection" {
	name = "tfAccSiteCollection<RANDOM_PLACEHOLDER>"
}
resource "luminate_collection" "new-collection" {
	name = "tfAccCollectionForApp<RANDOM_PLACEHOLDER>"
}
resource "luminate_collection_site_link" "new-collection-site-link" {
	site_id = "${luminate_site.new-site-collection.id}"
	collection_ids = sort(["${luminate_collection.new-collection.id}"])
}
resource "luminate_rdp_application" "new-rdp-application-collection" {
	site_id = "${luminate_site.new-site-collection.id}"
	collection_id = "${luminate_collection.new-collection.id}"
	name = "tfAccRDPWithCollection<RANDOM_PLACEHOLDER>"
	internal_address = "tcp://127.0.0.2"
    depends_on = [luminate_collection_site_link.new-collection-site-link]
}
`

func TestAccLuminateRDPApplication(t *testing.T) {
	resourceName := "luminate_rdp_application.new-rdp-application"
	resourceNameCollection := "luminate_rdp_application.new-rdp-application-collection"
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:  strings.ReplaceAll(testAccRDPApplication_minimal, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfAccRDP%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "visible", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "internal_address", "tcp://127.0.0.2:3389"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tfaccrdp%d.rdp.%s", randNum, testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tfaccrdp%d.rdp.%s", randNum, testAccDomain)),
				),
			},
			{
				Config:  strings.ReplaceAll(testAccRDPApplication_changeInternalAddress, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Destroy: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_address", "tcp://127.0.0.2:33"),
				),
			},
			{
				Config:  strings.ReplaceAll(testAccRDPApplication_changeInternalAddress_2, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Destroy: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_address", "tcp://127.0.0.2:3389"),
				),
			},
			{
				Config:  strings.ReplaceAll(testAccRDPApplication_changeInternalAddress_3, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Destroy: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "internal_address", "tcp://127.0.0.3:3389"),
				),
			},
			{
				Config: strings.ReplaceAll(testAccRDPApplication_options, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfAccRDPUpd%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "internal_address", "tcp://127.0.0.5:126"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tfaccrdp%d.rdp.%s", randNum, testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tfaccrdp%d.rdp.%s", randNum, testAccDomain)),
				),
			},
			{
				Config: strings.ReplaceAll(testAccRDPApplication_collection, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", fmt.Sprintf("tfAccRDPWithCollection%d", randNum))),
			},
		},
	})
}
