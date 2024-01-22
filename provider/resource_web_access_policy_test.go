package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceWebAccessPolicy_enabled = `
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
		enabled = "true"
		name =  "resourceWebAccessPolicy_enabled"
		identity_provider_id = "local"

		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430","ed974d59-1941-4584-9336-2a9ed35043f2"]
  		applications = ["${luminate_web_application.new-application.id}"]
		group_ids = ["3b61849d-f08d-42d3-a158-da1a53cd2ac6"]
	}
`

const resourceWebAccessPolicy_collection = `
	resource "luminate_site" "new-site-collection" {
	   name = "tfAccSiteAccessPolicyCollection"
	}
	resource "luminate_collection" "new-collection" {
		name = "tfAccCollectionForAppCollection"
	}
	resource "luminate_collection_site_link" "new-collection-site-link" {
		site_id = "${luminate_site.new-site-collection.id}"
		collection_ids = sort(["${luminate_collection.new-collection.id}"])
	}
	resource "luminate_web_application" "new-application-collection" {
	 site_id = "${luminate_site.new-site-collection.id}"
	 collection_id = "${luminate_collection.new-collection.id}"
	 name = "tfAccApplicationAccessPolicyCollection"
	 internal_address = "http://127.0.0.1:8080"
	 depends_on = [luminate_collection_site_link.new-collection-site-link]
	}
	resource "luminate_web_access_policy" "new-web-access-policy-collection" {
		enabled = "true"
		name =  "resourceWebAccessPolicy_collection"
	 	collection_id = "${luminate_collection.new-collection.id}"
		identity_provider_id = "local"

		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_web_application.new-application-collection.id}"]
	 	depends_on = [luminate_collection_site_link.new-collection-site-link]
	}
`

const resourceWebAccessPolicy_disabled = `
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
		enabled = "false"
  		name =  "resourceWebAccessPolicy_disabled"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_web_application.new-application.id}"]
	}
`

const resourceWebAccessPolicy_enabled_not_specified = `
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
  		name =  "resourceWebAccessPolicy_enabled_not_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_web_application.new-application.id}"]
	}
`

const resourceWebAccessPolicy_conditions_specified = `
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
  		name =  "resourceWebAccessPolicy_conditions_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_web_application.new-application.id}"]

		conditions {
    		source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
    		location = ["Wallis and Futuna"]

    		managed_device {
      			opswat = true
      			symantec_web_security_service = true
    		}
  		}

	}
`

const resourceWebAccessPolicy_conditions_specified_update = `
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
		name =  "resourceWebAccessPolicy_conditions_specified"
		identity_provider_id = "local"
	
		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
		applications = ["${luminate_web_application.new-application.id}"]
	
		conditions {
			source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
			location = ["Wallis and Futuna"]
	
			managed_device {
				opswat = false
				symantec_web_security_service = true
			}
		}

	}
`

const resourceWebAccessPolicy_validators_specified = `
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
  		name =  "resourceWebAccessPolicy_validators_specified"
		identity_provider_id = "local"
	
  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_web_application.new-application.id}"]
	}
`

func TestAccLuminateWebAccessPolicy(t *testing.T) {
	resourceName := "luminate_web_access_policy.new-web-access-policy"
	resourceNameCollection := "luminate_web_access_policy.new-web-access-policy-collection"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceWebAccessPolicy_enabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceWebAccessPolicy_enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: resourceWebAccessPolicy_disabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceWebAccessPolicy_disabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: resourceWebAccessPolicy_enabled_not_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceWebAccessPolicy_enabled_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config:  resourceWebAccessPolicy_conditions_specified,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceWebAccessPolicy_conditions_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.0", "127.0.0.1/24"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.1", "1.1.1.1/16"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.location.0", "Wallis and Futuna"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.managed_device.0.opswat", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.managed_device.0.symantec_web_security_service", "true"),
				),
			},
			{
				Config: resourceWebAccessPolicy_conditions_specified_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceWebAccessPolicy_conditions_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.0", "127.0.0.1/24"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.1", "1.1.1.1/16"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.location.0", "Wallis and Futuna"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.managed_device.0.opswat", "false"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.managed_device.0.symantec_web_security_service", "true"),
				),
			},
			{
				Config: resourceWebAccessPolicy_validators_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceWebAccessPolicy_validators_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: resourceWebAccessPolicy_collection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", "resourceWebAccessPolicy_collection"),
				),
			},
		},
	})
}
