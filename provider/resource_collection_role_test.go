package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const testCollectionRole = `
	resource "luminate_collection" "collection" {
		name = "collectionToBeAssign"
	} 
	resource "luminate_collection_role" "policy-owner" {
		role_type = "PolicyOwner"
		identity_provider_id =  "local"
		entity_id = "f75f45b8-d10d-4aa6-9200-5c6d60110430"
		entity_type = "User"
		collection_id = "${luminate_collection.collection.id}"
	}
	resource "luminate_collection_role" "app-owner" {
		role_type = "ApplicationOwner"
		identity_provider_id =  "local"
		entity_id = "f75f45b8-d10d-4aa6-9200-5c6d60110430"
		entity_type = "User"
		collection_id = "${luminate_collection.collection.id}"
	}
`

func TestAccLuminateCollectionRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCollectionRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("luminate_collection_role.policy-owner", "role_type", "PolicyOwner"),
					resource.TestCheckResourceAttr("luminate_collection_role.app-owner", "role_type", "ApplicationOwner"),
				),
			},
		},
	})
}
