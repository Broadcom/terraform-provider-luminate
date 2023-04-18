package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testAccLuminateTenantRoleBindings = `
	resource "luminate_tenant_role" "tenant-admins" {
		role = "TenantAdmin"
		entity_id = "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"
		identity_provider_id = "local"
}
`

func TestAccLuminateTenantRoleBindings(t *testing.T) {
	const resourceName = "luminate_tenant_role.tenant-admins"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:  testAccLuminateTenantRoleBindings,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role", "TenantAdmin"),
					resource.TestCheckResourceAttr(resourceName, "entity_id", "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"),
				),
			},
		},
	})
}
