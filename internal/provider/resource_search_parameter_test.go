package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSearchParameter_elementNameAndPatternFilterInExpression(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireLegacyMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearchParameter_elementNameAndPatternFilterInExpression,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "name", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "module", "fhir-4.0.1"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "type", "string"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "expression.0", "telecom"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "expression.1", "system|phone"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "expression.2", "value"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "reference.0.resource_id", "Patient"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "reference.0.resource_type", "Entity"),
				),
			},
		},
	})
}

func TestAccResourceSearchParameter_elementNameAndIndexInExpression(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireLegacyMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearchParameter_elementNameAndIndexInExpression,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "name", "first-name"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "module", ""),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "type", "string"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "expression.0", "name"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "expression.1", "0"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "expression.2", "given"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "expression.3", "0"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "reference.0.resource_id", "Patient"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "reference.0.resource_type", "Entity"),
				),
			},
		},
	})
}

func TestAccResourceSearchParameter_extension(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireLegacyMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearchParameter_extension,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "id", "Appointment.custom-date"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "name", "custom-date"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "type", "date"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "expression.0", "meta"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "expression.1", "extension"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "expression.2", "url|https://fhir.patientsknowbest.com/structuredefinition/custom-date"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "expression.3", "valueDateTime"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "reference.0.resource_id", "Appointment"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "reference.0.resource_type", "Entity"),
				),
			},
		},
	})
}

const testAccResourceSearchParameter_elementNameAndPatternFilterInExpression = `
resource "aidbox_search_parameter" "example_phone" {
  name     = "phone-number"
  module   = "fhir-4.0.1"
  type     = "string"
  reference {
    resource_id   = "Patient"
    resource_type = "Entity"
  }
  expression = ["telecom", "system|phone", "value"]
}
`

const testAccResourceSearchParameter_elementNameAndIndexInExpression = `
resource "aidbox_search_parameter" "example_name" {
  name     = "first-name"
  type     = "string"
  reference {
    resource_id   = "Patient"
    resource_type = "Entity"
  }
  expression = ["name", "0", "given", "0"]
}
`

const testAccResourceSearchParameter_extension = `
resource "aidbox_search_parameter" "example_extension" {
  name     = "custom-date"
  type     = "date"
  reference {
    resource_id   = "Appointment"
    resource_type = "Entity"
  }
  expression = ["meta", "extension", "url|https://fhir.patientsknowbest.com/structuredefinition/custom-date", "valueDateTime"]
}
`
