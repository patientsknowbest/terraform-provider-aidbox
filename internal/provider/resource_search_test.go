package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "param_parser", "reference"),
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
					resource.TestCheckResourceAttrPtr("aidbox_search.example_phone", "id", &previousIdState),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "name", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "param_parser", "reference"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "module", "fhir-4.0.1"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "reference.0.resource_id", "Patient"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "reference.0.resource_type", "Entity"),
					resource.TestCheckResourceAttr("aidbox_search.example_phone", "where", "phone-number = 00000000000"),
				),
			},
		},
	})
}

func TestAccResourceSearch_tokenSql(t *testing.T) {
	previousIdState := ""
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearch_tokenSql,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "name", "prescription-status"),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "param_parser", "token"),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "reference.0.resource_id", "ServiceRequest"),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "reference.0.resource_type", "Entity"),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "token_sql.0.only_code", "resource @> 'dispensed'"),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "token_sql.0.no_system", "resource @> 'dispensed'"),
					resource.TestCheckResourceAttrWith("aidbox_search.example_prescription_status", "id", func(id string) error {
						previousIdState = id
						return nil
					}),
				),
			},
			{
				Config: testAccResourceSearch_tokenSql_updateClause,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("aidbox_search.example_prescription_status", "id", &previousIdState),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "name", "prescription-status"),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "param_parser", "token"),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "reference.0.resource_id", "ServiceRequest"),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "reference.0.resource_type", "Entity"),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "token_sql.0.only_code", "resource @> 'fulfilled'"),
					resource.TestCheckResourceAttr("aidbox_search.example_prescription_status", "token_sql.0.no_system", "resource @> 'fulfilled'"),
				),
			},
		},
	})
}

const testAccResourceSearch_happyPath = `
resource "aidbox_search" "example_phone" {
  name         = "phone-number"
  param_parser = "reference"
  module       = "fhir-4.0.1"
  reference {
    resource_id   = "Patient"
    resource_type = "Entity"
  }
  where = "phone-number = 01125636365"
}
`

const testAccResourceSearch_happyPath_updateWhereClause = `
resource "aidbox_search" "example_phone" {
  name         = "phone-number"
  param_parser = "reference"
  module       = "fhir-4.0.1"
  reference {
    resource_id   = "Patient"
    resource_type = "Entity"
  }
  where = "phone-number = 00000000000"
}
`

const testAccResourceSearch_tokenSql = `
resource "aidbox_search" "example_prescription_status" {
  name         = "prescription-status"
  param_parser = "token"
  reference {
    resource_id   = "ServiceRequest"
    resource_type = "Entity"
  }
  where = "true"
  token_sql {
    only_code = "resource @> 'dispensed'"
    no_system = "resource @> 'dispensed'"
  }
}
`

const testAccResourceSearch_tokenSql_updateClause = `
resource "aidbox_search" "example_prescription_status" {
  name         = "prescription-status"
  param_parser = "token"
  reference {
    resource_id   = "ServiceRequest"
    resource_type = "Entity"
  }
  where = "true"
  token_sql {
    only_code = "resource @> 'fulfilled'"
    no_system = "resource @> 'fulfilled'"
  }
}
`
