package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceBox(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testMultiboxProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBox,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_box.mybox", "id", "mybox"),
					resource.TestCheckResourceAttr("aidbox_box.mybox", "fhir_version", "fhir-3.0.1"),
				),
			},
		},
	})
}

func TestAccResourceBox_insideResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testMultiboxProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBoxWithInternalResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_box.mybox", "id", "mybox"),
					resource.TestCheckResourceAttr("aidbox_box.mybox", "fhir_version", "fhir-3.0.1"),
					resource.TestCheckResourceAttr("aidbox_token_introspector.example", "type", "jwt"),
				),
			},
		},
	})
}

const testAccResourceBox = `
resource "aidbox_box" "mybox" {
  name = "mybox"
  fhir_version  = "fhir-3.0.1" 
  description = "A box instance within multibox, a multi-tenant aidbox server"
}
`

const testAccResourceBoxWithInternalResource = `
resource "aidbox_box" "mybox" {
  name = "mybox"
  fhir_version  = "fhir-3.0.1" 
  description = "A box instance within multibox, a multi-tenant aidbox server"
}
resource "aidbox_token_introspector" "example" {
  box_id = aidbox_box.mybox.name
  type = "jwt"
  jwks_uri = "http://keycloak:8080/auth/realms/pkb/protocol/openid-connect/certs"
  jwt {
    iss = "http://keycloak:8080/auth/realms/pkb"
  }
}
`
