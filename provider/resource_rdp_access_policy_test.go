package provider

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceRdpAccessPolicy_enabled = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
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
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.2"
	}
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
		enabled = "false"
  		name =  "resourceRdpAccessPolicy_disabled"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_rdp_application.new-rdp-application.id}"]

  		allow_long_term_password = "true"
	}
`

const resourceRdpAccessPolicy_WebRdp_default_settings = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.2"
		sub_type = "RDP_BROWSER_SINGLE_MACHINE"
	}
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
  		name =  "resourceRdpAccessPolicy_WebRdp_default_settings"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_rdp_application.new-rdp-application.id}"]

  		target_protocol_subtype = "RDP_BROWSER"
	}
`

const resourceRdpAccessPolicy_WebRdp_custom_settings = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.2"
		sub_type = "RDP_BROWSER_SINGLE_MACHINE"
	}
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
  		name =  "resourceRdpAccessPolicy_WebRdp_custom_settings"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_rdp_application.new-rdp-application.id}"]

  		target_protocol_subtype = "RDP_BROWSER"
		web_rdp_settings {
			disable_copy  = true
			disable_paste = true
		}
	}
`

const resourceRdpAccessPolicy_enabled_not_specified = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.2"
	}
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
  		name =  "resourceRdpAccessPolicy_enabled_not_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_rdp_application.new-rdp-application.id}"]
	}
`

const resourceRdpAccessPolicy_no_long_term_password_specified = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.2"
	}
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
		enabled = "true"
  		name =  "resourceRdpAccessPolicy_no_long_term_password_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_rdp_application.new-rdp-application.id}"]
	}
`

const resourceRdpAccessPolicy_conditions_specified = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.2"
	}
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
  		name =  "resourceRdpAccessPolicy_conditions_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_rdp_application.new-rdp-application.id}"]

		conditions {
    		source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
    		location = ["Wallis and Futuna"]
  		}

	}
`

const resourceRdpAccessPolicy_validators_specified = `
	resource "luminate_site" "new-site" {
		name = "tfAccSite<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_rdp_application" "new-rdp-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccRDP<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.2"
	}
	resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
  		name =  "resourceRdpAccessPolicy_validators_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_rdp_application.new-rdp-application.id}"]

		validators {
			web_verification = true
		}
	}
`
const resourceRdpAccessPolicy_collection = `
	resource "luminate_site" "new-site-collection" {
		name = "tfAccSiteCollection<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_collection" "new-collection" {
		name = "tfAccCollectionForAppCollection<RANDOM_PLACEHOLDER>"
	}
	resource "luminate_collection_site_link" "new-collection-site-link" {
		site_id = "${luminate_site.new-site-collection.id}"
		collection_ids = sort(["${luminate_collection.new-collection.id}"])
	}
	resource "luminate_rdp_application" "new-rdp-application-collection" {
		site_id = "${luminate_site.new-site-collection.id}"
		name = "tfAccRDPCollection<RANDOM_PLACEHOLDER>"
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
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: strings.ReplaceAll(resourceRdpAccessPolicy_enabled, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_long_term_password", "true"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceRdpAccessPolicy_disabled, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_disabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "allow_long_term_password", "true"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceRdpAccessPolicy_WebRdp_default_settings, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_WebRdp_default_settings"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "target_protocol_subtype", "RDP_BROWSER"),
					resource.TestCheckResourceAttr(resourceName, "web_rdp_settings.0.disable_copy", "false"),
					resource.TestCheckResourceAttr(resourceName, "web_rdp_settings.0.disable_paste", "true"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceRdpAccessPolicy_WebRdp_custom_settings, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_WebRdp_custom_settings"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "target_protocol_subtype", "RDP_BROWSER"),
					resource.TestCheckResourceAttr(resourceName, "web_rdp_settings.0.disable_copy", "true"),
					resource.TestCheckResourceAttr(resourceName, "web_rdp_settings.0.disable_paste", "true"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceRdpAccessPolicy_enabled_not_specified, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_enabled_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceRdpAccessPolicy_no_long_term_password_specified, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_no_long_term_password_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_long_term_password", "false"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceRdpAccessPolicy_conditions_specified, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_conditions_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.0", "127.0.0.1/24"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.1", "1.1.1.1/16"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.location.0", "Wallis and Futuna"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceRdpAccessPolicy_validators_specified, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceRdpAccessPolicy_validators_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "validators.0.web_verification", "true"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceRdpAccessPolicy_collection, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", "resourceRdpAccessPolicy_collection"),
				),
			},
		},
	})
}
