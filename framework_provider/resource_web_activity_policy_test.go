package framework_provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
		name =  "tfAccWebActivityPolicy_minimal%d"
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
		name =  "tfAccWebActivityPolicy_enabled%d"
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

func resourceWebActivityPolicy_collection(userID1 string, rand int) string {
	return fmt.Sprintf(`
	resource "luminate_site" "new-site-collection" {
	   name = "tfAccSiteActivityPolicyCollection%d"
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
	 name = "tfAccAppActivityPolicyCollection%d"
	 internal_address = "http://127.0.0.1:8080"
	 depends_on = [luminate_collection_site_link.new-collection-site-link]
	}
	resource "luminate_web_activity_policy" "new-web-activity-policy-collection" {
		enabled = "true"
		name =  "tfAccWebActivityPolicy_collection%d"
	 	collection_id = "${luminate_collection.new-collection.id}"
		identity_provider_id = "local"

		rules = [
			{
				action = "BLOCK"
				conditions = {
					file_uploaded = true
				}
			}
		]

		user_ids = ["%s"]
  		applications = ["${luminate_web_application.new-application-collection.id}"]
	 	depends_on = [luminate_collection_site_link.new-collection-site-link]
	}`, rand, rand, rand, rand, userID1)
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
  		name =  "tfAccWebActivityPolicy_disabled%d"
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
  		name =  "tfAccWebActivityPolicy_enabled_not_specified%d"
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
  		name =  "tfAccWebActivityPolicy_conditions_specified%d"
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
		name =  "tfAccWebActivityPolicy_conditions_specified_update%d"
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
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: resourceWebActivityPolicy_minimal(randNum),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("tfAccWebActivityPolicy_minimal%d", randNum)),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("BLOCK"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("file_uploaded"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				Config: resourceWebActivityPolicy_enabled(groupName, userID1, userID2, randNum),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("tfAccWebActivityPolicy_enabled%d", randNum)),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("BLOCK_USER"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("uri_accessed"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("arguments").AtMapKey("uri_list").AtSliceIndex(0),
						knownvalue.StringExact("/admin"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("arguments").AtMapKey("uri_list").AtSliceIndex(1),
						knownvalue.StringExact("/users"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(1).AtMapKey("action"),
						knownvalue.StringExact("DISCONNECT_USER"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(1).AtMapKey("conditions").AtMapKey("file_uploaded"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				Config: resourceWebActivityPolicy_disabled(userID1, randNum),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("tfAccWebActivityPolicy_disabled%d", randNum)),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("enabled"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("DISCONNECT_USER"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("http_command"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("arguments").AtMapKey("commands").AtSliceIndex(0),
						knownvalue.StringExact("GET"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("arguments").AtMapKey("commands").AtSliceIndex(1),
						knownvalue.StringExact("POST"),
					),
				},
			},
			{
				Config: resourceWebActivityPolicy_enabled_not_specified(userID1, randNum),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("tfAccWebActivityPolicy_enabled_not_specified%d", randNum)),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("BLOCK_USER"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("file_downloaded"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				Config:  resourceWebActivityPolicy_conditions_specified(userID1, randNum),
				Destroy: false,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("tfAccWebActivityPolicy_conditions_specified%d", randNum)),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("conditions").AtMapKey("source_ip").AtSliceIndex(0),
						knownvalue.StringExact("127.0.0.1/24"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("conditions").AtMapKey("source_ip").AtSliceIndex(1),
						knownvalue.StringExact("1.1.1.1/16"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("conditions").AtMapKey("location").AtSliceIndex(0),
						knownvalue.StringExact("Wallis and Futuna"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("conditions").AtMapKey("managed_device").AtMapKey("opswat"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("conditions").AtMapKey("managed_device").AtMapKey("symantec_web_security_service"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("BLOCK_USER"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("uri_accessed"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("arguments").AtMapKey("uri_list").AtSliceIndex(0),
						knownvalue.StringExact("/admin"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("arguments").AtMapKey("uri_list").AtSliceIndex(1),
						knownvalue.StringExact("/users"),
					),
				},
			},
			{
				Config:  resourceWebActivityPolicy_conditions_specified_update(userID1, randNum),
				Destroy: false,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("tfAccWebActivityPolicy_conditions_specified_update%d", randNum)),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("conditions").AtMapKey("source_ip").AtSliceIndex(0),
						knownvalue.StringExact("127.0.0.1/32"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("conditions").AtMapKey("location").AtSliceIndex(0),
						knownvalue.StringExact("Canada"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("conditions").AtMapKey("managed_device").AtMapKey("opswat"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("conditions").AtMapKey("unmanaged_device").AtMapKey("symantec_web_security_service"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("BLOCK_USER"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("uri_accessed"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("arguments").AtMapKey("uri_list").AtSliceIndex(0),
						knownvalue.StringExact("/admin"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("arguments").AtMapKey("uri_list").AtSliceIndex(1),
						knownvalue.StringExact("/users"),
					),
				},
			},
		},
	})
}

func TestAccLuminateResourceWebActivityPolicyWithCollection(t *testing.T) {
	resourceName := "luminate_web_activity_policy.new-web-activity-policy-collection"
	userID1, _, _ := getUsersAndGroupsFromEnvVars(t)
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: resourceWebActivityPolicy_collection(userID1, randNum),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("tfAccWebActivityPolicy_collection%d", randNum)),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("BLOCK"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("conditions").AtMapKey("file_uploaded"),
						knownvalue.Bool(true),
					),
				},
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
