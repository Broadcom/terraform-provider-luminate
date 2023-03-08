package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const resourceTcpAccessPolicy_enabled = `
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
		enabled = "true"
		name =  "resourceTcpAccessPolicy_enabled"
		identity_provider_id = "local"

		user_ids = ["e9bb7894-a6e4-44de-a2b3-ee9e5e72485a"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]

		allow_temporary_token = "true"
		allow_public_key = "true"
	}
`

const resourceTcpAccessPolicy_disabled = `
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
		enabled = "false"
  		name =  "resourceTcpAccessPolicy_disabled"
		identity_provider_id = "local"

  		user_ids = ["e9bb7894-a6e4-44de-a2b3-ee9e5e72485a"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]
	}
`

const resourceTcpAccessPolicy_enabled_not_specified = `
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
  		name =  "resourceTcpAccessPolicy_enabled_not_specified"
		identity_provider_id = "local"

  		user_ids = ["e9bb7894-a6e4-44de-a2b3-ee9e5e72485a"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]
	}
`

const resourceTcpAccessPolicy_optional_not_specified = `
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
		enabled = "true"
  		name =  "resourceTcpAccessPolicy_optional_not_specified"
		identity_provider_id = "local"

  		user_ids = ["e9bb7894-a6e4-44de-a2b3-ee9e5e72485a"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]
	}
`

const resourceTcpAccessPolicy_conditions_specified = `
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
  		name =  "resourceTcpAccessPolicy_conditions_specified"
		identity_provider_id = "local"

  		user_ids = ["e9bb7894-a6e4-44de-a2b3-ee9e5e72485a"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]

		conditions {
    		source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
    		location = ["Wallis and Futuna"]
  		}

	}
`

const resourceTcpAccessPolicy_validators_specified = `
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
  		name =  "resourceTcpAccessPolicy_validators_specified"
		identity_provider_id = "local"

  		user_ids = ["e9bb7894-a6e4-44de-a2b3-ee9e5e72485a"]
  		applications = ["7fdde321-c795-4a49-82e1-210ee9a8e1de"]

		validators {
			web_verification = true
		}
	}
`

func TestAccLuminateTcpAccessPolicy(t *testing.T) {
	resourceName := "luminate_tcp_access_policy.new-tcp-access-policy"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTcpAccessPolicy_enabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_temporary_token", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_public_key", "true"),
				),
			},
			{
				Config: resourceTcpAccessPolicy_disabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_disabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: resourceTcpAccessPolicy_enabled_not_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_enabled_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: resourceTcpAccessPolicy_optional_not_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_optional_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_temporary_token", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_public_key", "false"),
				),
			},
			{
				Config: resourceTcpAccessPolicy_conditions_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_conditions_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.0", "127.0.0.1/24"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.1", "1.1.1.1/16"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.location.0", "Wallis and Futuna"),
				),
			},
			{
				Config: resourceTcpAccessPolicy_validators_specified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_validators_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "validators.0.web_verification", "true"),
				),
			},
		},
	})
}
