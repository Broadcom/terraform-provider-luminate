package provider

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func resourceTcpAccessPolicy_enabled(rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicyTCP%d"
	}
	resource "luminate_tcp_application" "new-tcp-application" {
	  name = "tfAccTCPAccessPolicy%d"
	  site_id = "${luminate_site.new-site.id}"
	  target {
		address = "127.0.0.1"
		ports = ["8080"]
	  }
	}
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
		enabled = "true"
		name =  "resourceTcpAccessPolicy_enabled"
		identity_provider_id = "local"

		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_tcp_application.new-tcp-application.id}"]

		allow_temporary_token = "true"
		allow_public_key = "true"
	}
`, rand, rand)
}

func resourceTcpAccessPolicy_disabled(rand int) string {
	return fmt.Sprintf(
		`
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicyTCP%d"
	}
	resource "luminate_tcp_application" "new-tcp-application" {
	  name = "tfAccTCPAccessPolicy%d"
	  site_id = "${luminate_site.new-site.id}"
	  target {
		address = "127.0.0.1"
		ports = ["8080"]
		port_mapping = [80]	
	  }
	}
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
		enabled = "false"
  		name =  "resourceTcpAccessPolicy_disabled"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_tcp_application.new-tcp-application.id}"]
	} `, rand, rand)
}

func resourceTcpAccessPolicy_enabled_not_specified(rand int) string {
	return fmt.Sprintf(
		`
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicyTCP%d"
	}
	resource "luminate_tcp_application" "new-tcp-application" {
	  name = "tfAccTCPAccessPolicy%d"
	  site_id = "${luminate_site.new-site.id}"
	  target {
		address = "127.0.0.1"
		ports = ["8080"]
		port_mapping = [80]	
	  }
	}
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
  		name =  "resourceTcpAccessPolicy_enabled_not_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_tcp_application.new-tcp-application.id}"]
	}
`, rand, rand)
}

func resourceTcpAccessPolicy_optional_not_specified(rand int) string {
	return fmt.Sprintf(
		`
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicyTCP%d"
	}
	resource "luminate_tcp_application" "new-tcp-application" {
	  name = "tfAccTCPAccessPolicy%d"
	  site_id = "${luminate_site.new-site.id}"
	  target {
		address = "127.0.0.1"
		ports = ["8080"]
		port_mapping = [80]	
	  }
	}
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
		enabled = "true"
  		name =  "resourceTcpAccessPolicy_optional_not_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_tcp_application.new-tcp-application.id}"]
	}
`, rand, rand)
}

func resourceTcpAccessPolicy_conditions_specified(rand int) string {
	return fmt.Sprintf(
		`
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicyTCP%d"
	}
	resource "luminate_tcp_application" "new-tcp-application" {
	  name = "tfAccTCPAccessPolicy%d"
	  site_id = "${luminate_site.new-site.id}"
	  target {
		address = "127.0.0.1"
		ports = ["8080"]
		port_mapping = [80]	
	  }
	}
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
  		name =  "resourceTcpAccessPolicy_conditions_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_tcp_application.new-tcp-application.id}"]

		conditions {
    		source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
    		location = ["Wallis and Futuna"]
  		}

	}`, rand, rand)
}

func resourceTcpAccessPolicy_validators_specified(rand int) string {
	return fmt.Sprintf(
		`
	resource "luminate_site" "new-site" {
		name = "tfAccSiteAccessPolicyTCP%d"
	}
	resource "luminate_tcp_application" "new-tcp-application" {
	  name = "tfAccTCPAccessPolicy%d"
	  site_id = "${luminate_site.new-site.id}"
	  target {
		address = "127.0.0.1"
		ports = ["8080"]
		port_mapping = [80]	
	  }
	}
	resource "luminate_tcp_access_policy" "new-tcp-access-policy" {
  		name =  "resourceTcpAccessPolicy_validators_specified"
		identity_provider_id = "local"

  		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_tcp_application.new-tcp-application.id}"]

		validators {
			web_verification = true
		}
	}
`, rand, rand)
}

func resourceTCPAccessPolicy_collection(rand int) string {
	return fmt.Sprintf(
		`
	resource "luminate_site" "new-site-collection" {
	   name = "tfAccSiteAccessPolicyCollection%d"
	}
	resource "luminate_collection" "new-collection" {
		name = "tfAccCollectionForApp%d"
	}
	resource "luminate_collection_site_link" "new-collection-site-link" {
		site_id = "${luminate_site.new-site-collection.id}"
		collection_ids = sort(["${luminate_collection.new-collection.id}"])
	}
	resource "luminate_tcp_application" "new-tcp-application-collection" {
	  name = "tfAccTCPAccessPolicyCollection%d"
	  site_id = "${luminate_site.new-site-collection.id}"
      collection_id = "${luminate_collection.new-collection.id}"
	  target {
		address = "127.0.0.1"
		ports = ["8080"]
		port_mapping = [80]	
	  }

      depends_on = [luminate_collection_site_link.new-collection-site-link]
	}
	resource "luminate_tcp_access_policy" "new-tcp-access-policy-collection" {
		enabled = "true"
		name =  "resourceTcpAccessPolicy_collection"
      	collection_id = "${luminate_collection.new-collection.id}"
		identity_provider_id = "local"

		user_ids = ["f75f45b8-d10d-4aa6-9200-5c6d60110430"]
  		applications = ["${luminate_tcp_application.new-tcp-application-collection.id}"]


		allow_temporary_token = "true"
		allow_public_key = "true"

		depends_on = [luminate_collection_site_link.new-collection-site-link]
	}
`, rand, rand, rand)
}

func TestAccLuminateTcpAccessPolicy(t *testing.T) {
	resourceName := "luminate_tcp_access_policy.new-tcp-access-policy"
	resourceNameCollection := "luminate_tcp_access_policy.new-tcp-access-policy-collection"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: resourceTcpAccessPolicy_enabled(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_temporary_token", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_public_key", "true"),
				),
			},
			{
				Config: resourceTcpAccessPolicy_disabled(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_disabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: resourceTcpAccessPolicy_enabled_not_specified(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_enabled_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: resourceTcpAccessPolicy_optional_not_specified(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_optional_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_temporary_token", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_public_key", "false"),
				),
			},
			{
				Config: resourceTcpAccessPolicy_conditions_specified(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_conditions_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.0", "127.0.0.1/24"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.source_ip.1", "1.1.1.1/16"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.location.0", "Wallis and Futuna"),
				),
			},
			{
				Config: resourceTcpAccessPolicy_validators_specified(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceTcpAccessPolicy_validators_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "validators.0.web_verification", "true"),
				),
			},
			{
				Config: resourceTCPAccessPolicy_collection(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", "resourceTcpAccessPolicy_collection"),
				),
			},
		},
	})
}
