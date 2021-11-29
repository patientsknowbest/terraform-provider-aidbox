package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTokenIntrospector(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTokenIntrospector,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"TokenIntrospector_resource.foo", "sample_attribute", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccResourceTokenIntrospector = `
resource "aidbox_token_introspector" "example" {
  jwks_uri = "http://keycloak:8080/auth/realms/pkb/protocol/openid-connect/certs"
  jwt {
    iss = "http://keycloak:8080/auth/realms/pkb"
  }
}
`
