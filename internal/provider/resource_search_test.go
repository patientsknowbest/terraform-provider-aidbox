package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccResourceSearch_happyPath(t *testing.T) {
	previousIdState := ""
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearch_happyPath,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "name", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "module", "fhir-4.0.1"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "reference.0.resource_id", "Patient"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "reference.0.resource_type", "Entity"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "where", "phone-number = 01125636365"),
					resource.TestCheckResourceAttrWith("aidbox_search.example_phone", "id", func(id string) error {
						previousIdState = id
						return nil
					}),
				),
			},
			{
				Config: testAccResourceSearch_happyPath_updateWhereClause,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "name", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "module", "fhir-4.0.1"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "reference.0.resource_id", "Patient"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "reference.0.resource_type", "Entity"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "where", "phone-number = 00000000000"),
					resource.TestCheckResourceAttrWith("aidbox_search.example_phone", "id", func(id string) error {
						assert.Equalf(t, previousIdState, id, "Resource logical id unexpectedly changed after resource update")
						return nil
					}),
				),
			},
		},
	})
}

const testAccResourceSearch_happyPath = `
resource "aidbox_search" "example_phone" {
  name     = "phone-number"
  module   = "fhir-4.0.1"
  reference {
    resource_id   = "Patient"
    resource_type = "Entity"
  }
  where = "phone-number = 01125636365"
}
`

const testAccResourceSearch_happyPath_updateWhereClause = `
resource "aidbox_search" "example_phone" {
  name     = "phone-number"
  module   = "fhir-4.0.1"
  reference {
    resource_id   = "Patient"
    resource_type = "Entity"
  }
  where = "phone-number = 00000000000"
}
`
