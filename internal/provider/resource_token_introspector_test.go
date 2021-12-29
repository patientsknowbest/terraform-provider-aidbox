package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTokenIntrospector(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTokenIntrospector,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_token_introspector.example", "jwks_uri", "http://keycloak:8080/auth/realms/pkb/protocol/openid-connect/certs"),
					resource.TestCheckResourceAttr("aidbox_token_introspector.example", "jwt.0.iss", "http://keycloak:8080/auth/realms/pkb"),
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
