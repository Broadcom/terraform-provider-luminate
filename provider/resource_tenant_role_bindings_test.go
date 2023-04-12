package provider

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testAccLuminateRoleBindings = `
	resource "lumiante_tenant_role" "tenant-admins" {
	  role = "TenantAdmin"
	  entity = [{id:"24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9", "identity_provider_id":"local" }]
	}
`

func TestAccLuminateRoleBindings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLuminateRoleBindings,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lumiante_tenant_role.tenant-admins", "role", "TenantAdmin"),
					resource.TestCheckResourceAttr("lumiante_tenant_role.tenant-admins", "entity.0.id", "24d8dcf9-b95c-4c92-a1a6-21083eb4d3a9"),
					resource.TestCheckResourceAttr("lumiante_tenant_role.tenant-admins", "entity.0.identity_provider_id", "local"),
				),
			},
		},
	})
}
