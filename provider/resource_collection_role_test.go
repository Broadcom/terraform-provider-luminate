package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func testCollectionRole(name string) string {
	return fmt.Sprintf(`
	data "luminate_user"  "my-users" {
		identity_provider_id = "local"
		users = ["%s"]
	}
	resource "luminate_collection" "collection" {
		name = "collectionToBeAssign"
	} 
	resource "luminate_collection_role" "policy-owner" {
		role_type = "PolicyOwner"
		identity_provider_id =  "local"
		entity_id = "${data.luminate_user.my-users.user_ids.0}"
		entity_type = "User"
		collection_id = "${luminate_collection.collection.id}"
	}
	resource "luminate_collection_role" "app-owner" {
		role_type = "ApplicationOwner"
		identity_provider_id =  "local"
		entity_id = "${data.luminate_user.my-users.user_ids.0}"
		entity_type = "User"
		collection_id = "${luminate_collection.collection.id}"
	}`, name)
}

func TestAccLuminateCollectionRole(t *testing.T) {
	var username string
	if username = os.Getenv("TEST_USERNAME"); username == "" {
		t.Error("skipping TestAccLuminateDataSourceUser no username provided")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCollectionRole(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("luminate_collection_role.policy-owner", "role_type", "PolicyOwner"),
					resource.TestCheckResourceAttr("luminate_collection_role.app-owner", "role_type", "ApplicationOwner"),
				),
			},
		},
	})
}
