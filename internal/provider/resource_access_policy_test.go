package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAccessPolicy_schema(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAccessPolicy_schema,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_access_policy.example", "engine", "json-schema"),
				),
			},
		},
	})
}

const testAccResourceAccessPolicy_schema = `
resource "aidbox_access_policy" "example" {
  description = "A policy to allow postman to access data"
  engine = "json-schema"
  # The test cries about whitespace differences after application when using jsonencode. 
  # For practical purposes, jsonencode is much easier to read, so you should use that.
  schema = "{\"required\":[\"client\",\"uri\",\"request-method\"],\"properties\":{\"uri\":{\"type\":\"string\",\"pattern\":\"^/fhir/.*\"},\"client\":{\"required\":[\"id\"],\"properties\":{\"id\":{\"const\":\"postman\"}}},\"request-method\":{\"const\":\"get\"}}}"
  #schema = jsonencode({
  #  "required" = [
  #    "client",
  #    "uri",
  #    "request-method" ]
  #  "properties" = {
  #    "uri" = {
  #      "type" = "string"
  #      "pattern" = "^/fhir/.*"
  #    }
  #    "client" = {
  #      "required" = ["id"]
  #      "properties" = {
  #        "id" = {
  #          const = "postman"
  #        }
  #      }
  #    }
  #    "request-method" = {
  #      "const" = "get"
  #    }
  #  }
  #})
}
`
