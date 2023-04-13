package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testAccLuminateRoleBindings = `
	resource "luminate_tenant_role" "tenant-admins" {
		role = "TenantAdmin"
		entity_id = "b28cef1c-ced4-441e-9ecc-5887cabcda60"
		identity_provider_id = "local"
}
`

func TestAccLuminateRoleBindings(t *testing.T) {
	const resourceName = "luminate_tenant_role.tenant-admins"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:  testAccLuminateRoleBindings,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role", "TenantAdmin"),
					resource.TestCheckResourceAttr(resourceName, "entity_id", "b28cef1c-ced4-441e-9ecc-5887cabcda60"),
				),
			},
		},
	})
}
