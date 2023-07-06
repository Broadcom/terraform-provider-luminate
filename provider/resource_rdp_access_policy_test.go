package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const resourceRdpAccessPolicy_enabled = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP"
		internal_address = "tcp://127.0.0.2"
	}
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
		enabled = "true"
  		name =  "resourceRdpAccessPolicy_enabled"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_rdp_application.new-rdp-application.id}"]

  		allow_long_term_password = "true"
	}
`

const resourceRdpAccessPolicy_disabled = `
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
		enabled = "false"
  		name =  "resourceRdpAccessPolicy_disabled"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]

  		allow_long_term_password = "true"
	}
`

const resourceRdpAccessPolicy_enabled_not_specified = `
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
  		name =  "resourceRdpAccessPolicy_enabled_not_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]

	}
`

const resourceRdpAccessPolicy_no_long_term_password_specified = `
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
		enabled = "true"
  		name =  "resourceRdpAccessPolicy_no_long_term_password_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]
	}
`

const resourceRdpAccessPolicy_conditions_specified = `
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
  		name =  "resourceRdpAccessPolicy_conditions_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]

		conditions {
    		source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
    		location = ["Wallis and Futuna"]
  		}

	}
`

const resourceRdpAccessPolicy_validators_specified = `
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
  		name =  "resourceRdpAccessPolicy_validators_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]

		validators {
			web_verification = true
		}
	}
`
const resourceRdpAccessPolicy_collection = `
	resource "luminate_site" "new-site-collection" {
		name = "tfAccSiteCollection"
	}
	resource "luminate_collection" "new-collection" {
		name = "tfAccCollectionForAppCollection"
	}
	resource "luminate_collection_site_link" "new-collection-site-link" {
		site_id = "${luminate_site.new-site-collection.id}"
		collection_ids = sort(["${luminate_collection.new-collection.id}"])
	}
	resource "luminate_rdp_application" "new-rdp-application-collection" {
		site_id = "${luminate_site.new-site-collection.id}"
		name = "tfAccRDPCollection"
		internal_address = "tcp://127.0.0.2"
        depends_on = [luminate_collection_site_link.new-collection-site-link]
      	collection_id = "${luminate_collection.new-collection.id}"
	}
	resource "luminate_rdp_access_policy" "new-rdp-access-policy-collection" {
		enabled = "true"
  		name =  "resourceRdpAccessPolicy_collection"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_rdp_application.new-rdp-application-collection.id}"]

  		allow_long_term_password = "true"
        depends_on = [luminate_collection_site_link.new-collection-site-link]
	    collection_id = "${luminate_collection.new-collection.id}"
		
	}
`

func TestAccLuminateRdpAccessPolicy(t *testing.T) {
	resourceName := "luminate_rdp_access_policy.new-rdp-access-policy"
	resourceNameCollection := "luminate_rdp_access_policy.new-rdp-access-policy-collection"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceRdpAccessPolicy_enabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_long_term_password", "true"),
				),
			},
			{
				Config: resourceRdpAccessPolicy_disabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_disabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "allow_long_term_password", "true"),
				),
			},
			{
				Config: resourceRdpAccessPolicy_enabled_not_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_enabled_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: resourceRdpAccessPolicy_no_long_term_password_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_no_long_term_password_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_long_term_password", "false"),
				),
			},
			{
				Config: resourceRdpAccessPolicy_conditions_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_conditions_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.0", "127.0.0.1/24"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.1", "1.1.1.1/16"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.location.0", "Wallis and Futuna"),
				),
			},
			{
				Config: resourceRdpAccessPolicy_validators_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_validators_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "validators.0.web_verification", "true"),
				),
			},
			{
				Config: resourceRdpAccessPolicy_collection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", "resourceRdpAccessPolicy_collection"),
				),
			},
		},
	})
}
