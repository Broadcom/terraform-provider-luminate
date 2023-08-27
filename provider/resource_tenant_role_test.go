package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const testTenantRole = `
	resource "luminate_tenant_role" "tenant-admin" {
		role_type = "TenantAdmin"
		identity_provider_id =  "local"
		entity_id = "f75f45b8-d10d-4aa6-9200-5c6d60110430"
		entity_type = "User"
	}
	resource "luminate_tenant_role" "tenant-viewer" {
		role_type = "TenantViewer"
		identity_provider_id =  "local"
		entity_id = "f75f45b8-d10d-4aa6-9200-5c6d60110430"
		entity_type = "User"
	}
`

func TestAccLuminateTenantRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTenantRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("luminate_tenant_role.tenant-admin", "role_type", "TenantAdmin"),
					resource.TestCheckResourceAttr("luminate_tenant_role.tenant-viewer", "role_type", "TenantViewer"),
				),
			},
		},
	})
}
