package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"math/rand"
	"testing"
	"time"
)

func TestAccLuminateCollectionRoleBindings(t *testing.T) {
	const resourceNameApp = "luminate_collection_role.app-owner"
	const resourceNamePolicy = "luminate_collection_role.policy-owner"
	const collectionName = "luminate_collection.new-collection"
	rand.Seed(time.Now().UnixNano())
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLuminateCollectionRoleBindingsAppOwner("tfAccCollectionRoleApp", randNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameApp, "role", "ApplicationOwner"),
					resource.TestCheckResourceAttr(resourceNameApp, "entity_id", "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"),
					resource.TestCheckResourceAttrPair(resourceNameApp, "collection_id", collectionName, "id"),
				),
			},
			{
				Config: testAccLuminateCollectionRoleBindingsPolicyOwner("tfAccCollectionRolePolicy", randNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePolicy, "role", "PolicyOwner"),
					resource.TestCheckResourceAttr(resourceNamePolicy, "entity_id", "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"),
					resource.TestCheckResourceAttrPair(resourceNamePolicy, "collection_id", collectionName, "id"),
				),
			},
		},
	})
}

func testAccLuminateCollectionRoleBindingsAppOwner(collectionName string, rand int) string {
	return fmt.Sprintf(`resource "luminate_collection" "new-collection" {
		name = "%s%d"
	}
	resource "luminate_collection_role" "app-owner" {
		role = "ApplicationOwner"
		entity_id = "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"
		identity_provider_id = "local"
		collection_id = "${luminate_collection.new-collection.id}"
	}
`, collectionName, rand)
}

func testAccLuminateCollectionRoleBindingsPolicyOwner(collectionName string, rand int) string {
	return fmt.Sprintf(`resource "luminate_collection" "new-collection" {
		name = "%s%d"
	}
	resource "luminate_collection_role" "policy-owner" {
		role = "PolicyOwner"
		entity_id = "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"
		identity_provider_id = "local"
		collection_id = "${luminate_collection.new-collection.id}"
	}
`, collectionName, rand)
}
