package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTokenIntrospector_jwt(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTokenIntrospector_jwt,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_token_introspector.example", "type", "jwt"),
					resource.TestCheckResourceAttr("aidbox_token_introspector.example", "jwks_uri", "http://keycloak:8080/auth/realms/pkb/protocol/openid-connect/certs"),
					resource.TestCheckResourceAttr("aidbox_token_introspector.example", "jwt.0.iss", "http://keycloak:8080/auth/realms/pkb"),
				),
			},
		},
	})
}

func TestAccResourceTokenIntrospector_opaque(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTokenIntrospector_opaque,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_token_introspector.example2", "type", "opaque"),
					resource.TestCheckResourceAttr("aidbox_token_introspector.example2", "introspection_endpoint.0.authorization", "Bearer foobar"),
					resource.TestCheckResourceAttr("aidbox_token_introspector.example2", "introspection_endpoint.0.url", "https://example.com/auth"),
				),
			},
		},
	})
}

const testAccResourceTokenIntrospector_jwt = `
resource "aidbox_token_introspector" "example" {
  type = "jwt"
  jwks_uri = "http://keycloak:8080/auth/realms/pkb/protocol/openid-connect/certs"
  jwt {
    iss = "http://keycloak:8080/auth/realms/pkb"
  }
}
`

const testAccResourceTokenIntrospector_opaque = `
resource "aidbox_token_introspector" "example2" {
  type = "opaque"
  introspection_endpoint {
    authorization = "Bearer foobar"
    url = "https://example.com/auth"
  }
}
`
