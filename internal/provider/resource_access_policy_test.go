package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
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
	previousIdState := ""
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAccessPolicy_schema,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("aidbox_access_policy.example", "schema", func(valueFromServer string) error {
						assert.True(t, jsonDiffSuppressFunc("", schema_v1, valueFromServer, nil), "Value received from server does not match (semantically): %s", valueFromServer)
						return nil
					}),
					resource.TestCheckResourceAttrWith("aidbox_access_policy.example", "id", func(id string) error {
						previousIdState = id
						return nil
					}),
				),
			},
			{
				Config: testAccResourceAccessPolicy_schema_updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("aidbox_access_policy.example", "schema", func(valueFromServer string) error {
						assert.True(t, jsonDiffSuppressFunc("", schema_v2, valueFromServer, nil), "Value received from server does not match (semantically): %s", valueFromServer)
						return nil
					}),
					resource.TestCheckResourceAttrWith("aidbox_access_policy.example", "id", func(id string) error {
						assert.Equalf(t, previousIdState, id, "Resource logical id unexpectedly changed after resource update")
						return nil
					}),
				),
			},
		},
	})
}

const testAccResourceAccessPolicy_schema = `
resource "aidbox_access_policy" "example" {
  description = "A policy to allow postman to access data"
  engine = "json-schema"
  schema = <<-EOT` + schema_v1 + `
EOT
}
`

const schema_v1 = `
{
    "properties":
    {
        "client":
        {
            "properties":
            {
                "id":
                {
                    "const": "postman"
                }
            },
            "required":
            [
                "id"
            ]
        },
        "request-method":
        {
            "const": "get"
        },
        "uri":
        {
            "pattern": "^/fhir/.*",
            "type": "string"
        }
    },
    "required":
    [
        "client",
        "uri",
        "request-method"
    ]
}`

const testAccResourceAccessPolicy_schema_updated = `
resource "aidbox_access_policy" "example" {
  description = "A policy to allow postman to access data"
  engine = "json-schema"
  schema = <<-EOT` + schema_v2 + `
EOT
}
`

const schema_v2 = `
{
    "properties":
    {
        "client":
        {
            "properties":
            {
                "id":
                {
                    "const": "postman"
                }
            },
            "required":
            [
                "id"
            ]
        },
        "request-method":
        {
            "const": "get"
        },
        "uri":
        {
            "pattern": "^/(fhir|ValueSet)(/.*|$)",
            "type": "string"
        }
    },
    "required":
    [
        "client",
        "uri",
        "request-method"
    ]
}`

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
