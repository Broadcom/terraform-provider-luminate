package provider

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceSshAccessPolicy_enabled = `
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicySSH<RANDOM_PLACEHOLDER>"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.5"
	}
	resource "luminate_ssh_access_policy" "new-ssh-access-policy" {
		enabled = "true"
		name =  "resourceSshAccessPolicy_enabled<RANDOM_PLACEHOLDER>"
		identity_provider_id = "local"

		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_ssh_application.new-ssh-application.id}"]

		accounts = ["ubuntu", "ec2-user"]
		use_auto_mapping = "true"
		allow_agent_forwarding = "true"
		allow_temporary_token = "true"
		allow_public_key = "true"
	}
`

const resourceSshAccessPolicy_disabled = `
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicySSH<RANDOM_PLACEHOLDER>"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.5"
	}
	resource "luminate_ssh_access_policy" "new-ssh-access-policy" {
		enabled = "false"
  		name =  "resourceSshAccessPolicy_disabled<RANDOM_PLACEHOLDER>"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_ssh_application.new-ssh-application.id}"]

		accounts = ["ubuntu", "ec2-user"]
		allow_temporary_token = "true"

	}
`

const resourceSshAccessPolicy_enabled_not_specified = `
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicySSH<RANDOM_PLACEHOLDER>"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.5"
	}
	resource "luminate_ssh_access_policy" "new-ssh-access-policy" {
  		name =  "resourceSshAccessPolicy_enabled_not_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_ssh_application.new-ssh-application.id}"]

		accounts = ["ubuntu", "ec2-user"]
		allow_temporary_token = "true"

	}
`

const resourceSshAccessPolicy_optional_not_specified = `
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicySSH<RANDOM_PLACEHOLDER>"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.5"
	}
	resource "luminate_ssh_access_policy" "new-ssh-access-policy" {
		enabled = "true"
  		name =  "resourceSshAccessPolicy_optional_not_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_ssh_application.new-ssh-application.id}"]

		accounts = ["ubuntu", "ec2-user"]
	}
`

const resourceSshAccessPolicy_conditions_specified = `
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicySSH<RANDOM_PLACEHOLDER>"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.5"
	}
	resource "luminate_ssh_access_policy" "new-ssh-access-policy" {
  		name =  "resourceSshAccessPolicy_conditions_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_ssh_application.new-ssh-application.id}"]

		conditions {
    		source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
    		location = ["Wallis and Futuna"]
  		}

		accounts = ["ubuntu", "ec2-user"]
		allow_temporary_token = "true"
	}
`

const resourceSshAccessPolicy_validators_specified = `
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicySSH<RANDOM_PLACEHOLDER>"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd<RANDOM_PLACEHOLDER>"
		internal_address = "tcp://127.0.0.5"
	}
	resource "luminate_ssh_access_policy" "new-ssh-access-policy" {
  		name =  "resourceSshAccessPolicy_validators_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_ssh_application.new-ssh-application.id}"]

		validators {
			web_verification = true
		}
		allow_temporary_token = "true"

		accounts = ["ubuntu", "ec2-user"]
	}
`
const resourceSshAccessPolicy_Collection = `
	resource "luminate_site" "new-site-collection" {
		name = "tfAccSiteAccessPolicySSHCollection"
	}
	resource "luminate_collection" "new-collection" {
		name = "tfAccCollectionForAppCollection"
	}
	resource "luminate_collection_site_link" "new-collection-site-link" {
		site_id = "${luminate_site.new-site-collection.id}"
		collection_ids = sort(["${luminate_collection.new-collection.id}"])
	}
	resource "luminate_ssh_application" "new-ssh-application-collection" {
		site_id = "${luminate_site.new-site-collection.id}"
		name = "tfAccSSHCollection<RANDOM_PLACEHOLDER>"
      	collection_id = "${luminate_collection.new-collection.id}"
		internal_address = "tcp://127.0.0.5"
        depends_on = [luminate_collection_site_link.new-collection-site-link]
	}
	resource "luminate_ssh_access_policy" "new-ssh-access-policy-collection" {
		enabled = "true"
		name =  "resourceSshAccessPolicy_Collection"
		identity_provider_id = "local"
      	collection_id = "${luminate_collection.new-collection.id}"

		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_ssh_application.new-ssh-application-collection.id}"]

		accounts = ["ubuntu", "ec2-user"]
		use_auto_mapping = "true"
		allow_agent_forwarding = "true"
		allow_temporary_token = "true"
		allow_public_key = "true"
		depends_on = [luminate_collection_site_link.new-collection-site-link]
	}
`

func TestAccLuminateSshAccessPolicy(t *testing.T) {
	resourceName := "luminate_ssh_access_policy.new-ssh-access-policy"
	resourceNameCollection := "luminate_ssh_access_policy.new-ssh-access-policy-collection"
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: strings.ReplaceAll(resourceSshAccessPolicy_enabled, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("resourceSshAccessPolicy_enabled%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1", "ec2-user"),
					resource.TestCheckResourceAttr(resourceName, "use_auto_mapping", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_agent_forwarding", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_temporary_token", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_public_key", "true"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceSshAccessPolicy_disabled, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("resourceSshAccessPolicy_disabled%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1", "ec2-user"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceSshAccessPolicy_enabled_not_specified, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceSshAccessPolicy_enabled_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1", "ec2-user"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceSshAccessPolicy_optional_not_specified, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceSshAccessPolicy_optional_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1", "ec2-user"),
					resource.TestCheckResourceAttr(resourceName, "allow_temporary_token", "true"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceSshAccessPolicy_conditions_specified, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceSshAccessPolicy_conditions_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1", "ec2-user"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.0", "127.0.0.1/24"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.1", "1.1.1.1/16"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.location.0", "Wallis and Futuna"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceSshAccessPolicy_validators_specified, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceSshAccessPolicy_validators_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1", "ec2-user"),
					resource.TestCheckResourceAttr(resourceName, "validators.0.web_verification", "true"),
				),
			},
			{
				Config: strings.ReplaceAll(resourceSshAccessPolicy_Collection, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", "resourceSshAccessPolicy_Collection"),
				),
			},
		},
	})
}
