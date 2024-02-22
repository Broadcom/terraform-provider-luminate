package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func testTenantRole(name string) string {
	return fmt.Sprintf(`
	resource "luminate_tenant_role" "tenant-admin" {
		role_type = "TenantAdmin"
		identity_provider_id =  "local"
		entity_id = "%s"
		entity_type = "User"
	}
	resource "luminate_tenant_role" "tenant-viewer" {
		role_type = "TenantViewer"
		identity_provider_id =  "local"
		entity_id = "%s"
		entity_type = "User"
	}`, name, name)
}

func TestAccLuminateTenantRole(t *testing.T) {
	var userID string
	if userID = os.Getenv("TEST_USER_ID"); userID == "" {
		t.Skip("skipping TestAccLuminateTenantRole no user id provided")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTenantRole(userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("luminate_tenant_role.tenant-admin", "role_type", "TenantAdmin"),
					resource.TestCheckResourceAttr("luminate_tenant_role.tenant-viewer", "role_type", "TenantViewer"),
				),
			},
		},
	})
}
