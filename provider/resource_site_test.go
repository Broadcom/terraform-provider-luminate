package provider

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourceSite_minimal = `
resource "luminate_site" "new-site" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
}
`

func getTestSiteResourceConfigAndCheckFunc(region string, authenticationMode dto.SiteAuthenticationMode) (string, resource.TestCheckFunc) {
	randNum := 100 + rand.Intn(100)
	resourceName, config := testAccResourceSite_options(region, randNum, authenticationMode)
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfAccSiteOpt%d", randNum)),
		resource.TestCheckResourceAttr(resourceName, "region", region),
		resource.TestCheckResourceAttr(resourceName, "authentication_mode", string(authenticationMode)),
		resource.TestCheckResourceAttr(resourceName, "mute_health_notification", "true"),
		resource.TestCheckResourceAttr(resourceName, "kubernetes_persistent_volume_name", "K8SVolume"),
	)

	return config, check
}

func testAccResourceSite_options(region string, randNum int, authenticationMode dto.SiteAuthenticationMode) (string, string) {
	resourceName := fmt.Sprintf("new_site%d", randNum)
	fullResourceName := fmt.Sprintf("luminate_site.%s", resourceName)
	return fullResourceName, fmt.Sprintf(`resource "luminate_site" "%s" {
	name = "tfAccSiteOpt%d"
	region = "%s"
  	authentication_mode = "%s"
	mute_health_notification = "true"
	kubernetes_persistent_volume_name = "K8SVolume"
}`, resourceName, randNum, region, authenticationMode)
}

func TestAccLuminateSite(t *testing.T) {
	resourceName := "luminate_site.new-site"
	var region string
	if region = os.Getenv("TEST_SITE_REGION"); region == "" {
		t.Error("stopping TestAccLuminateSite no  site provided")
	}
	randNum := 100 + rand.Intn(100)

	configWithConnectorMode, checkWithConnectorMode := getTestSiteResourceConfigAndCheckFunc(region, dto.SiteAuthenticationModeManual)
	configWithSiteMode, checkWithSiteMode := getTestSiteResourceConfigAndCheckFunc(region, dto.SiteAuthenticationModeOrchestrator)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtocol6Providers,
		Steps: []resource.TestStep{
			{
				Config: strings.ReplaceAll(testAccResourceSite_minimal, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tfAccSite%d", randNum)),
					resource.TestCheckResourceAttr(resourceName, "mute_health_notification", "false"),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_persistent_volume_name", ""),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", dto.SiteAuthenticationModeManual),
				),
			},
			{
				Config: configWithConnectorMode,
				Check:  checkWithConnectorMode,
			},
			{
				Config: configWithSiteMode,
				Check:  checkWithSiteMode,
			},
		},
	})
}
