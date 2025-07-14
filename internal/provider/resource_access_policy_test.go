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

func TestAccResourceAccessPolicy_allow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAccessPolicy_allow,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_access_policy.mypolicy", "engine", "allow"),
				),
			},
		},
	})
}

func TestAccResourceAccessPolicy_schema_updated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAccessPolicy_schema,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_access_policy.example", "schema", "{\"required\":[\"client\",\"uri\",\"request-method\"],\"properties\":{\"uri\":{\"type\":\"string\",\"pattern\":\"^/fhir/.*\"},\"client\":{\"required\":[\"id\"],\"properties\":{\"id\":{\"const\":\"postman\"}}},\"request-method\":{\"const\":\"get\"}}}"),
				),
			},
			{
				Config: testAccResourceAccessPolicy_schema_updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_access_policy.example", "schema", "{\"required\":[\"client\",\"uri\",\"request-method\"],\"properties\":{\"uri\":{\"type\":\"string\",\"pattern\":\"^/(fhir|ValueSet)(/.*|$)\"},\"client\":{\"required\":[\"id\"],\"properties\":{\"id\":{\"const\":\"postman\"}}},\"request-method\":{\"const\":\"get\"}}}"),
				),
			},
		},
	})
}

const testAccResourceAccessPolicy_schema = `
resource "aidbox_access_policy" "example" {
  description = "A policy to allow postman to access data"
  engine = "json-schema"
  # The test complains about whitespace differences after application when using jsonencode. 
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

const testAccResourceAccessPolicy_schema_updated = `
resource "aidbox_access_policy" "example" {
  description = "A policy to allow postman to access data"
  engine = "json-schema"
  # The test complains about whitespace differences after application when using jsonencode. 
  # For practical purposes, jsonencode is much easier to read, so you should use that.
  schema = "{\"required\":[\"client\",\"uri\",\"request-method\"],\"properties\":{\"uri\":{\"type\":\"string\",\"pattern\":\"^/(fhir|ValueSet)(/.*|$)\"},\"client\":{\"required\":[\"id\"],\"properties\":{\"id\":{\"const\":\"postman\"}}},\"request-method\":{\"const\":\"get\"}}}"
  #schema = jsonencode({
  #  "required" = [
  #    "client",
  #    "uri",
  #    "request-method" ]
  #  "properties" = {
  #    "uri" = {
  #      "type" = "string"
  #      "pattern" = "^/(fhir|ValueSet)(/.*|$)"
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

const testAccResourceAccessPolicy_allow = `
resource "aidbox_client" "client" {
  name        = "client-id"
  secret      = "__sha256:2BB80D537B1DA3E38BD30361AA855686BDE0EACD7162FEF6A25FE97BF527A25B"
  grant_types = ["basic"]
}
resource "aidbox_access_policy" "mypolicy" {
  description = "A policy to allow client to access data"
  engine      = "allow"
  link {
    resource_id   = aidbox_client.client.name
    resource_type = "Client"
  }
}
`
