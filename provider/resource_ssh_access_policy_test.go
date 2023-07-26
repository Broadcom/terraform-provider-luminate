package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceSshAccessPolicy_enabled = `
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicySSH"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd"
		internal_address = "tcp://127.0.0.5"
	}
	resource "luminate_ssh_access_policy" "new-ssh-access-policy" {
		enabled = "true"
		name =  "resourceSshAccessPolicy_enabled"
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
		name = "tfAccSiteAccessPolicySSH"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd"
		internal_address = "tcp://127.0.0.5"
	}
	resource "luminate_ssh_access_policy" "new-ssh-access-policy" {
		enabled = "false"
  		name =  "resourceSshAccessPolicy_disabled"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_ssh_application.new-ssh-application.id}"]

		accounts = ["ubuntu", "ec2-user"]
		allow_temporary_token = "true"

	}
`

const resourceSshAccessPolicy_enabled_not_specified = `
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicySSH"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd"
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
		name = "tfAccSiteAccessPolicySSH"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd"
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
		name = "tfAccSiteAccessPolicySSH"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd"
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
		name = "tfAccSiteAccessPolicySSH"
	}
	
	resource "luminate_ssh_application" "new-ssh-application" {
		site_id = "${luminate_site.new-site.id}"
		name = "tfAccSSHUpd"
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
		name = "tfAccSSHCollection"
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

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceSshAccessPolicy_enabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceSshAccessPolicy_enabled"),
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
				Config: resourceSshAccessPolicy_disabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceSshAccessPolicy_disabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1", "ec2-user"),
				),
			},
			{
				Config: resourceSshAccessPolicy_enabled_not_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceSshAccessPolicy_enabled_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1", "ec2-user"),
				),
			},
			{
				Config: resourceSshAccessPolicy_optional_not_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceSshAccessPolicy_optional_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1", "ec2-user"),
					resource.TestCheckResourceAttr(resourceName, "allow_temporary_token", "true"),
				),
			},
			{
				Config: resourceSshAccessPolicy_conditions_specified,
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
				Config: resourceSshAccessPolicy_validators_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceSshAccessPolicy_validators_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "accounts.0", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "accounts.1", "ec2-user"),
					resource.TestCheckResourceAttr(resourceName, "validators.0.web_verification", "true"),
				),
			},
			{
				Config: resourceSshAccessPolicy_Collection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", "resourceSshAccessPolicy_Collection"),
				),
			},
		},
	})
}
