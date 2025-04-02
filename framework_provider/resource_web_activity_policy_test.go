package framework_provider

import (
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"math/rand"
	"os"
	"testing"
)

func resourceWebActivityPolicy_minimal(rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteActivityPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationActivityPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_activity_policy" "new-web-activity-policy" {
		name =  "resourceWebActivityPolicy_minimal%d"
		applications = ["${luminate_web_application.new-application.id}"]

		rules = [
			{
				action = "BLOCK"
				conditions = {
					file_uploaded = true
				}
			}
		]
	}`, rand, rand, rand)
}

func resourceWebActivityPolicy_enabled(groupName,
	userID1,
	userID2 string,
	rand int) string {
	return fmt.Sprintf(`
	data "luminate_group"  "my-groups" {
		identity_provider_id = "local"
		groups = ["%s"]
	}
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteActivityPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationActivityPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_activity_policy" "new-web-activity-policy" {
		enabled = "true"
		name =  "resourceWebActivityPolicy_enabled%d"
		identity_provider_id = "local"

		rules = [
			{
				action = "BLOCK_USER"
				conditions = {
					uri_accessed = true
					arguments = {
						uri_list = ["/admin", "/users"]
					}
				}
			},
			{
				action = "DISCONNECT_USER"
				conditions = {
					file_uploaded = true
				}
			}
		]

		user_ids = ["%s","%s"]
  		applications = ["${luminate_web_application.new-application.id}"]
		group_ids = ["${data.luminate_group.my-groups.group_ids.0}"]
	}`, groupName, rand, rand, rand, userID1, userID2)
}

func resourceWebActivityPolicy_disabled(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteActivityPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationActivityPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_activity_policy" "new-web-activity-policy" {
		enabled = "false"
  		name =  "resourceWebActivityPolicy_disabled%d"
		identity_provider_id = "local"

		rules = [
			{
				action = "DISCONNECT_USER"
				conditions = {
					http_command = true
					arguments = {
						commands = ["GET", "POST"]
					}
				}
			}
		]

  		user_ids = ["%s"]
  		applications = ["${luminate_web_application.new-application.id}"]
	}`, rand, rand, rand, userID1)
}

func resourceWebActivityPolicy_enabled_not_specified(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteActivityPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationActivityPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_activity_policy" "new-web-activity-policy" {
  		name =  "resourceWebActivityPolicy_enabled_not_specified%d"
		identity_provider_id = "local"

		rules = [
			{
				action = "BLOCK_USER"
				conditions = {
					file_downloaded = true
				}
			}
		]

  		user_ids = ["%s"]
  		applications = ["${luminate_web_application.new-application.id}"]
	}`, rand, rand, rand, userID1)
}

func resourceWebActivityPolicy_conditions_specified(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteActivityPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationActivityPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_activity_policy" "new-web-activity-policy" {
  		name =  "resourceWebActivityPolicy_conditions_specified%d"
		identity_provider_id = "local"

  		user_ids = ["%s"]
  		applications = ["${luminate_web_application.new-application.id}"]

		conditions = {
    		source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
    		location = ["Wallis and Futuna"]

    		managed_device = {
      			opswat = true
      			symantec_web_security_service = true
    		}
  		}

		rules = [
			{
				action = "BLOCK_USER"
				conditions = {
					uri_accessed = true
					arguments = {
						uri_list = ["/admin", "/users"]
					}
				}
			}
		]

	}
`, rand, rand, rand, userID1)
}

func resourceWebActivityPolicy_conditions_specified_update(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site" {
	   name = "tfAccSiteActivityPolicy%d"
	}
	resource "luminate_web_application" "new-application" {
	 site_id = "${luminate_site.new-site.id}"
	 name = "tfAccApplicationActivityPolicy%d"
	 internal_address = "http://127.0.0.1:8080"
	}
	resource "luminate_web_activity_policy" "new-web-activity-policy" {
		name =  "resourceWebActivityPolicy_conditions_specified%d"
		identity_provider_id = "local"
	
		user_ids = ["%s"]
		applications = ["${luminate_web_application.new-application.id}"]
	
		conditions = {
			source_ip = ["127.0.0.1/32"]
			location = ["Canada"]
	
			managed_device = {
				opswat = false
				symantec_web_security_service = true
			}

			unmanaged_device = {
				symantec_web_security_service = true
			}
		}

		rules = [
			{
				action = "BLOCK_USER"
				conditions = {
					uri_accessed = true
					arguments = {
						uri_list = ["/admin", "/users"]
					}
				}
			}
		]

	}`, rand, rand, rand, userID1)
}

func TestAccLuminateResourceWebActivityPolicyConditionsSpecifiedWithUpdate(t *testing.T) {
	resourceName := "luminate_web_activity_policy.new-web-activity-policy"
	userID1, userID2, groupName := getUsersAndGroupsFromEnvVars(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: resourceWebActivityPolicy_minimal(100 + rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", utils.CreateRegExpForNamePrefix("resourceWebActivityPolicy_minimal")),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "BLOCK"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.file_uploaded", "true"),
				),
			},
			{
				Config: resourceWebActivityPolicy_enabled(groupName, userID1, userID2, 100+rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", utils.CreateRegExpForNamePrefix("resourceWebActivityPolicy_enabled")),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "BLOCK_USER"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.uri_accessed", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.uri_list.0", "/admin"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.uri_list.1", "/users"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.action", "DISCONNECT_USER"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.conditions.file_uploaded", "true"),
				),
			},
			{
				Config: resourceWebActivityPolicy_disabled(userID1, 100+rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", utils.CreateRegExpForNamePrefix("resourceWebActivityPolicy_disabled")),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "DISCONNECT_USER"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.http_command", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.commands.0", "GET"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.commands.1", "POST"),
				),
			},
			{
				Config: resourceWebActivityPolicy_enabled_not_specified(userID1, 100+rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", utils.CreateRegExpForNamePrefix("resourceWebActivityPolicy_enabled_not_specified")),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "BLOCK_USER"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.file_downloaded", "true"),
				),
			},
			{
				Config:  resourceWebActivityPolicy_conditions_specified(userID1, 100+rand.Intn(100)),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", utils.CreateRegExpForNamePrefix("resourceWebActivityPolicy_conditions_specified")),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.source_ip.0", "127.0.0.1/24"),
					resource.TestCheckResourceAttr(resourceName, "conditions.source_ip.1", "1.1.1.1/16"),
					resource.TestCheckResourceAttr(resourceName, "conditions.location.0", "Wallis and Futuna"),
					resource.TestCheckResourceAttr(resourceName, "conditions.managed_device.opswat", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.managed_device.symantec_web_security_service", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "BLOCK_USER"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.uri_accessed", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.uri_list.0", "/admin"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.uri_list.1", "/users"),
				),
			},
			{
				Config: resourceWebActivityPolicy_conditions_specified_update(userID1, 100+rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", utils.CreateRegExpForNamePrefix("resourceWebActivityPolicy_conditions_specified")),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.source_ip.0", "127.0.0.1/32"),
					resource.TestCheckResourceAttr(resourceName, "conditions.location.0", "Canada"),
					resource.TestCheckResourceAttr(resourceName, "conditions.managed_device.opswat", "false"),
					resource.TestCheckResourceAttr(resourceName, "conditions.managed_device.symantec_web_security_service", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.unmanaged_device.symantec_web_security_service", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "BLOCK_USER"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.uri_accessed", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.uri_list.0", "/admin"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.uri_list.1", "/users"),
				),
			},
		},
	})
}

func getUsersAndGroupsFromEnvVars(t *testing.T) (string, string, string) {
	var userID1 string
	if userID1 = os.Getenv("TEST_USER_ID"); userID1 == "" {
		t.Error("stopping TestAccLuminateWebActivityPolicy no user id provided")
	}
	var userID2 string
	if userID2 = os.Getenv("TEST_USER_ID2"); userID2 == "" {
		t.Error("stopping TestAccLuminateWebActivityPolicy no user id 2 provided")
	}
	var groupName string
	if groupName = os.Getenv("TEST_GROUP_NAME"); groupName == "" {
		t.Error("stopping TestAccLuminateDataSourceGroup no group name provided")
	}
	return userID1, userID2, groupName
}
