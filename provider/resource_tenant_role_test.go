package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testTenantRole = `
	resource "luminate_tenant_role" "tenant-role" {
		role_type = "TenantAdmin"
		identity_provider_id =  "local"
		entity_id = "a8a48219-835f-4183-a2a9-bbba8cad8eb8"
		entity_type = "User"
	}
`

func TestAccLuminateTenantRole(t *testing.T) {
	resourceName := "luminate_tenant_role.tenant-role"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTenantRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_type", "TenantAdmin")),
			},
		},
	})
}
