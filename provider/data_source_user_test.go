package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceUser(name string) string {
	return fmt.Sprintf(
		`data "luminate_user"  "my-users" {
			identity_provider_id = "local"
			users = ["%s"]
		}`, name)
}

func TestAccLuminateDataSourceUser(t *testing.T) {
	resourceName := "data.luminate_user.my-users"
	var username, userID string

	if username = os.Getenv("TEST_USERNAME"); username == "" {
		t.Error("stopping TestAccLuminateDataSourceUser no  username provided")
	}
	if userID = os.Getenv("TEST_USER_ID"); userID == "" {
		t.Error("stopping TestAccLuminateDataSourceUser no user provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: newTestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUser(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_ids.0", userID),
				),
			},
		},
	})
}
