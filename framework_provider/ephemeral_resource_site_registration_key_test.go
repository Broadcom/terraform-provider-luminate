package framework_provider

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Broadcom/terraform-provider-luminate/test_utils"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const testAccEphemeralResourceSiteRegistrationKeyTemplate = `
resource "luminate_site" "new_site_<RANDOM_PLACEHOLDER>" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
  	authentication_mode = "orchestrator"
}

resource "luminate_site_registration_key_version" "site_registration_key_version" {
  version = <VERSION_PLACEHOLDER>
}

ephemeral "luminate_site_registration_key" "new_site_registration_key" {
	is_applying = terraform.applying
	site_id = luminate_site.new_site_<RANDOM_PLACEHOLDER>.id
	revoke_existing_key_immediately = true
	should_rotate = luminate_site_registration_key_version.site_registration_key_version.version_changed
}

provider "echo" {
  data = ephemeral.luminate_site_registration_key.new_site_registration_key
}

resource "echo" "<ECHO_NAME_PLACEHOLDER>" {}
`

const testAccEphemeralResourceSiteRegistrationKeyWithoutVersionTemplate = `
resource "luminate_site" "new_site_<RANDOM_PLACEHOLDER>" {
	name = "tfAccSite<RANDOM_PLACEHOLDER>"
	authentication_mode = "orchestrator"
}

ephemeral "luminate_site_registration_key" "new_site_registration_key" {
	is_applying                     = terraform.applying
	site_id                         = luminate_site.new_site_<RANDOM_PLACEHOLDER>.id
	revoke_existing_key_immediately = true
	should_rotate                   = <SHOULD_ROTATE_PLACEHOLDER>
}

provider "echo" {
	data = ephemeral.luminate_site_registration_key.new_site_registration_key
}

resource "echo" "<ECHO_NAME_PLACEHOLDER>" {}
`

func testAccEphemeralSiteRegistrationKeyProviders() map[string]func() (tfprotov6.ProviderServer, error) {
	providers := make(map[string]func() (tfprotov6.ProviderServer, error))
	for k, v := range testAccProtocol6Providers {
		providers[k] = v
	}
	providers["echo"] = echoprovider.NewProviderServer()
	return providers
}

func TestAccLuminateSiteRegistrationKeyWithKeyVersion(t *testing.T) {
	randNum := test_utils.GetRandomNumber()
	replacePlaceholders := func(config string, version string, echoName string) string {
		out := strings.ReplaceAll(config, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum))
		out = strings.ReplaceAll(out, "<VERSION_PLACEHOLDER>", version)
		out = strings.ReplaceAll(out, "<ECHO_NAME_PLACEHOLDER>", echoName)
		return out
	}

	replaceSiteName := func(template string) string {
		return strings.ReplaceAll(template, "tfAccSite", "tfAccSiteChange")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccEphemeralSiteRegistrationKeyProviders(),
		Steps: []resource.TestStep{
			{
				// 1) Apply with new version -> token should be set (not null)
				Config: replacePlaceholders(testAccEphemeralResourceSiteRegistrationKeyTemplate, "1", "token_apply_new"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.token_apply_new",
						tfjsonpath.New("data").AtMapKey("token"),
						knownvalue.StringRegexp(regexp.MustCompile(`.+`)),
					),
				},
			},
			{
				// 2) Plan phase -> computed token should be unknown (ephemeral returns null at plan). PreApply runs after plan, before apply.
				Config: replacePlaceholders(testAccEphemeralResourceSiteRegistrationKeyTemplate, "1", "token_plan"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// During plan the ephemeral returns null, so echo's data is unknown (no nested path in plan JSON)
						plancheck.ExpectUnknownValue(
							"echo.token_plan",
							tfjsonpath.New("data"),
						),
					},
				},
			},
			{
				// 3) Apply with version changed -> token should be set (not null)
				Config: replacePlaceholders(testAccEphemeralResourceSiteRegistrationKeyTemplate, "2", "token_apply_changed"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.token_apply_changed",
						tfjsonpath.New("data").AtMapKey("token"),
						knownvalue.StringRegexp(regexp.MustCompile(`.+`)),
					),
				},
			},
			{
				// 4) Apply again with version unchanged -> token should be null
				Config: replacePlaceholders(testAccEphemeralResourceSiteRegistrationKeyTemplate, "2", "token_apply_unchanged"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.token_apply_unchanged",
						tfjsonpath.New("data").AtMapKey("token"),
						knownvalue.StringExact(""),
					),
				},
			},
			{
				// 5) Plan only with version unchanged but site name changed -> token should be null
				Config:             replaceSiteName(replacePlaceholders(testAccEphemeralResourceSiteRegistrationKeyTemplate, "2", "token_apply_unchanged")),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.token_apply_unchanged",
						tfjsonpath.New("data").AtMapKey("token"),
						knownvalue.StringExact(""),
					),
				},
			},
			{
				// 6) Apply again with version unchanged but site name changed -> token should be null
				Config: replaceSiteName(replacePlaceholders(testAccEphemeralResourceSiteRegistrationKeyTemplate, "2", "token_apply_unchanged")),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.token_apply_unchanged",
						tfjsonpath.New("data").AtMapKey("token"),
						knownvalue.StringExact(""),
					),
				},
			},
			{
				// 7) Apply with version changed -> token should be set (not null)
				Config: replaceSiteName(replacePlaceholders(testAccEphemeralResourceSiteRegistrationKeyTemplate, "3", "token_apply_changed")),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.token_apply_changed",
						tfjsonpath.New("data").AtMapKey("token"),
						knownvalue.StringRegexp(regexp.MustCompile(`.+`)),
					),
				},
			},
		},
	})
}

func TestAccLuminateSiteRegistrationKeyWithoutKeyVersion(t *testing.T) {
	randNum := test_utils.GetRandomNumber()
	replacePlaceholders := func(config string, shouldRotate string, echoName string) string {
		out := strings.ReplaceAll(config, "<RANDOM_PLACEHOLDER>", strconv.Itoa(randNum))
		out = strings.ReplaceAll(out, "<SHOULD_ROTATE_PLACEHOLDER>", shouldRotate)
		out = strings.ReplaceAll(out, "<ECHO_NAME_PLACEHOLDER>", echoName)
		return out
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccEphemeralSiteRegistrationKeyProviders(),
		Steps: []resource.TestStep{
			{
				// 1) Apply with should_rotate = true -> token should be set (not null)
				Config: replacePlaceholders(testAccEphemeralResourceSiteRegistrationKeyWithoutVersionTemplate, "true", "token_rotate"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.token_rotate",
						tfjsonpath.New("data").AtMapKey("token"),
						knownvalue.StringRegexp(regexp.MustCompile(`.+`)),
					),
				},
			},
			{
				// 2) Apply with should_rotate = false -> token should be empty (new echo so it receives current ephemeral output)
				Config: replacePlaceholders(testAccEphemeralResourceSiteRegistrationKeyWithoutVersionTemplate, "false", "token_no_rotate"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.token_no_rotate",
						tfjsonpath.New("data").AtMapKey("token"),
						knownvalue.StringExact(""),
					),
				},
			},
			{
				// 3) Apply with should_rotate = true again -> token should be set (not null)
				Config: replacePlaceholders(testAccEphemeralResourceSiteRegistrationKeyWithoutVersionTemplate, "true", "token_rotate_again"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.token_rotate_again",
						tfjsonpath.New("data").AtMapKey("token"),
						knownvalue.StringRegexp(regexp.MustCompile(`.+`)),
					),
				},
			},
		},
	})
}
