package framework_provider

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const testAccResourceSite_minimal = `
resource "luminate_site2" "new-site" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
}
`

func testAccResourceSite_options(region string, randNum int) string {
	return fmt.Sprintf(`resource "luminate_site2" "new-site" {
	name = "tfAccSiteOpt%d"
	region = "%s"
	mute_health_notification = "true"
	kubernetes_persistent_volume_name = "K8SVolume"
}`, randNum, region)
}

func TestAccLuminateSite2(t *testing.T) {
	resourceName := "luminate_site2.new-site"
	var region string
	if region = os.Getenv("TEST_SITE_REGION"); region == "" {
		t.Error("stopping TestAccLuminateSite no  site provided")
	}
	randNum := 100 + rand.Intn(100)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: strings.ReplaceAll(testAccResourceSite_minimal, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("tfAccSite%d", randNum)),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("mute_health_notification"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("kubernetes_persistent_volume_name"),
						knownvalue.StringExact(""),
					),
				},
			},
			{
				Config: testAccResourceSite_options(region, randNum),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("tfAccSiteOpt%d", randNum)),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("region"),
						knownvalue.StringExact(region),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("mute_health_notification"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("kubernetes_persistent_volume_name"),
						knownvalue.StringExact("K8SVolume"),
					),
				},
			},
		},
	})
}
