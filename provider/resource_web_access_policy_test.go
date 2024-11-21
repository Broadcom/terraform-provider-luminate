package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"math/rand"
	"os"
	"regexp"
	"testing"
)

func resourceWebAccessPolicy_enabled(groupName,
	userID1,
	userID2 string,
	rand int) string {
	return fmt.Sprintf(`
	data "luminate_group"  "my-groups" {
		identity_provider_id = "local"
		groups = ["%s"]
	}
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
		enabled = "true"
		name =  "resourceWebAccessPolicy_enabled%d"
		identity_provider_id = "local"

		user_ids = ["%s","%s"]
  		applications = ["${luminate_web_application.new-application.id}"]
		group_ids = ["${data.luminate_group.my-groups.group_ids.0}"]
	}`, groupName, rand, rand, rand, userID1, userID2)
}

func resourceWebAccessPolicy_collection(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site-collection" {
	   name = "tfAccSiteAccessPolicyCollection%d"
	}
	resource "luminate_collection" "new-collection" {
		name = "tfAccCollectionForAppCollection%d"
	}
	resource "luminate_collection_site_link" "new-collection-site-link" {
		site_id = "${luminate_site.new-site-collection.id}"
		collection_ids = sort(["${luminate_collection.new-collection.id}"])
	}
	resource "luminate_web_application" "new-application-collection" {
	 site_id = "${luminate_site.new-site-collection.id}"
	 collection_id = "${luminate_collection.new-collection.id}"
	 name = "tfAccAppAccessPolicyCollection%d"
	 internal_address = "http://127.0.0.1:8080"
	 depends_on = [luminate_collection_site_link.new-collection-site-link]
	}
	resource "luminate_web_access_policy" "new-web-access-policy-collection" {
		enabled = "true"
		name =  "resourceWebAccessPolicy_collection%d"
	 	collection_id = "${luminate_collection.new-collection.id}"
		identity_provider_id = "local"

		user_ids = ["%s"]
  		applications = ["${luminate_web_application.new-application-collection.id}"]
	 	depends_on = [luminate_collection_site_link.new-collection-site-link]
	}`, rand, rand, rand, rand, userID1)
}

func resourceWebAccessPolicy_disabled(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
		enabled = "false"
  		name =  "resourceWebAccessPolicy_disabled%d"
		identity_provider_id = "local"

  		user_ids = ["%s"]
  		applications = ["${luminate_web_application.new-application.id}"]
	}`, rand, rand, rand, userID1)
}

func resourceWebAccessPolicy_enabled_not_specified(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
  		name =  "resourceWebAccessPolicy_enabled_not_specified%d"
		identity_provider_id = "local"

  		user_ids = ["%s"]
  		applications = ["${luminate_web_application.new-application.id}"]
	}`, rand, rand, rand, userID1)
}

func resourceWebAccessPolicy_conditions_specified(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
  		name =  "resourceWebAccessPolicy_conditions_specified%d"
		identity_provider_id = "local"

  		user_ids = ["%s"]
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
`, rand, rand, rand, userID1)
}

func resourceWebAccessPolicy_conditions_specified_update(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
		name =  "resourceWebAccessPolicy_conditions_specified%d"
		identity_provider_id = "local"
	
		user_ids = ["%s"]
		applications = ["${luminate_web_application.new-application.id}"]
	
		conditions {
			source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
			location = ["Wallis and Futuna"]
	
			managed_device {
				opswat = false
				symantec_web_security_service = true
			}
		}

	}`, rand, rand, rand, userID1)
}

func resourceWebAccessPolicy_validators_specified(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteAccessPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationAccessPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_access_policy" "new-web-access-policy" {
  		name =  "resourceWebAccessPolicy_validators_specified%d"
		identity_provider_id = "local"
	
  		user_ids = ["%s"]
  		applications = ["${luminate_web_application.new-application.id}"]
		validators {
			mfa = true
		}
	}`, rand, rand, rand, userID1)
}

func TestAccLuminateWebAccessPolicy(t *testing.T) {
	resourceName := "luminate_web_access_policy.new-web-access-policy"
	resourceNameCollection := "luminate_web_access_policy.new-web-access-policy-collection"
	var userID1 string
	if userID1 = os.Getenv("TEST_USER_ID"); userID1 == "" {
		t.Error("stopping TestAccLuminateWebAccessPolicy no user id provided")
	}
	var userID2 string
	if userID2 = os.Getenv("TEST_USER_ID2"); userID2 == "" {
		t.Error("stopping TestAccLuminateWebAccessPolicy no user id 2 provided")
	}
	var groupName string
	if groupName = os.Getenv("TEST_GROUP_NAME"); groupName == "" {
		t.Error("stopping TestAccLuminateDataSourceGroup no group name provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceWebAccessPolicy_enabled(groupName, userID1, userID2, 100+rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", createRegExpForNamePrefix("resourceWebAccessPolicy_enabled")),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: resourceWebAccessPolicy_disabled(userID1, 100+rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceWebAccessPolicy_disabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: resourceWebAccessPolicy_enabled_not_specified(userID1, 100+rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceWebAccessPolicy_enabled_not_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config:  resourceWebAccessPolicy_conditions_specified(userID1, 100+rand.Intn(100)),
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
				Config: resourceWebAccessPolicy_conditions_specified_update(userID1, 100+rand.Intn(100)),
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
				Config: resourceWebAccessPolicy_validators_specified(userID1, 100+rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "resourceWebAccessPolicy_validators_specified"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "validators.0.mfa", "true"),
				),
			},
			{
				Config: resourceWebAccessPolicy_collection(userID1, 100+rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCollection, "name", "resourceWebAccessPolicy_collection"),
				),
			},
		},
	})
}

func createRegExpForNamePrefix(prefix string) *regexp.Regexp {
	exp := fmt.Sprintf("^%s", prefix)
	return regexp.MustCompile(exp)
}
