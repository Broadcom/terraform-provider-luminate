package framework_provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"math/rand"
	"os"
	"regexp"
	"testing"
)

func resourceWebActivityPolicy_conditions_rules_specified(groupName,
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
		name =  "resourceWebActivityPolicy_conditions_rules_specified%d"
		identity_provider_id = "local"

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
				action = "DISCONNECT_USER"
				conditions = {
					file_downloaded = true
					file_uploaded = true
					uri_accessed = true
					http_command = true
					arguments = {
						uri_list = ["/admin", "/users"]
						commands = ["GET", "POST"]
					}
				}
			}
		]

		user_ids = ["%s","%s"]
  		applications = ["${luminate_web_application.new-application.id}"]
		group_ids = ["${data.luminate_group.my-groups.group_ids.0}"]
	}`, groupName, rand, rand, rand, userID1, userID2)
}

func TestAccLuminateResourceWebActivityPolicy(t *testing.T) {
	resourceName := "luminate_web_activity_policy.new-web-activity-policy"
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: resourceWebActivityPolicy_conditions_rules_specified(groupName, userID1, userID2, 100+rand.Intn(100)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", createRegExpForNamePrefix("resourceWebActivityPolicy_conditions_rules_specified")),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.source_ip.0", "127.0.0.1/24"),
					resource.TestCheckResourceAttr(resourceName, "conditions.source_ip.1", "1.1.1.1/16"),
					resource.TestCheckResourceAttr(resourceName, "conditions.location.0", "Wallis and Futuna"),
					resource.TestCheckResourceAttr(resourceName, "conditions.managed_device.opswat", "true"),
					resource.TestCheckResourceAttr(resourceName, "conditions.managed_device.symantec_web_security_service", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "DISCONNECT_USER"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.file_downloaded", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.file_uploaded", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.uri_accessed", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.http_command", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.uri_list.0", "/admin"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.uri_list.1", "/users"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.commands.0", "GET"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.arguments.commands.1", "POST"),
				),
			},
		},
	})
}

func createRegExpForNamePrefix(prefix string) *regexp.Regexp {
	exp := fmt.Sprintf("^%s", prefix)
	return regexp.MustCompile(exp)
}
