package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const terraformCode = `
data "luminate_group" "my-groups" {
  identity_provider_id = "e2f74d20-a1fa-4d26-b06c-b39244694788"
  groups               = ["styr-group"]
}

resource "luminate_rdp_access_policy" "rdp_access" {
  name                 = "test-rdp-policy"
  identity_provider_id = data.luminate_group.my-groups.identity_provider_id
  applications         = ["c7c49a80-1160-4600-9685-dff8b7e3d295"]
  group_ids            = data.luminate_group.my-groups.group_ids
}
`

func TestBugTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: terraformCode,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}
