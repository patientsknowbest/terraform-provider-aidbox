package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSearchParameterV2_elementNameAndPatternFilterInExpression(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearchParameterV2_elementNameAndPatternFilterInExpression,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "name", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "type", "string"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "base.0", "Patient"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "code", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "expression", "Patient.telecom.where(system = 'phone')"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "description", "Search patients by phone number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "url", "https://fhir.yourcompany.com/searchparameter/phone-number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "status", "active"),
				),
			},
		},
	})
}

func TestAccResourceSearchParameterV2_elementNameAndIndexInExpression(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearchParameterV2_elementNameAndIndexInExpression,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_name", "name", "first-name"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_name", "type", "string"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_name", "base.0", "Patient"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_name", "code", "first-name"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_name", "expression", "Patient.name.given"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_name", "description", "Search patients by first name"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_name", "url", "https://fhir.yourcompany.com/searchparameter/first-name"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_name", "status", "draft"),
				),
			},
		},
	})
}

func TestAccResourceSearchParameterV2_extension(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearchParameterV2_extension,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_extension", "name", "custom-date"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_extension", "type", "date"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_extension", "base.0", "Appointment"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_extension", "code", "custom-date"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_extension", "expression", "Appointment.meta.extension.where(url = 'https://fhir.yourcompany.com/structuredefinition/custom-date').valueDateTime"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_extension", "description", "Search appointments by custom date expression"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_extension", "url", "https://fhir.yourcompany.com/searchparameter/custom-date"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_extension", "status", "active"),
				),
			},
		},
	})
}

func TestAccResourceSearchParameterV2_update(t *testing.T) {
	previousIdState := ""
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearchParameterV2_elementNameAndPatternFilterInExpression,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "name", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "type", "string"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "base.0", "Patient"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "code", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "expression", "Patient.telecom.where(system = 'phone')"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "description", "Search patients by phone number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "url", "https://fhir.yourcompany.com/searchparameter/phone-number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "status", "active"),
					resource.TestCheckResourceAttrWith("aidbox_fhir_search_parameter.example_phone", "id", func(id string) error {
						previousIdState = id
						return nil
					}),
				),
			},
			{
				Config: testAccResourceSearchParameterV2_updatePhoneNumber,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("aidbox_fhir_search_parameter.example_phone", "id", &previousIdState),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "name", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "type", "string"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "base.0", "Patient"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "code", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "expression", "Patient.telecom.where(system = 'phone')"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "description", "Search patients by phone number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "url", "https://fhir.newdomain.com/searchparameter/phone-number"),
					resource.TestCheckResourceAttr("aidbox_fhir_search_parameter.example_phone", "status", "active"),
				),
			},
		},
	})
}

const testAccResourceSearchParameterV2_elementNameAndPatternFilterInExpression = `
resource "aidbox_fhir_search_parameter" "example_phone" {
  name        = "phone-number"
  type        = "string"
  base        = ["Patient"]
  code        = "phone-number"
  expression  = "Patient.telecom.where(system = 'phone')"
  description = "Search patients by phone number"
  url         = "https://fhir.yourcompany.com/searchparameter/phone-number"
}
`

const testAccResourceSearchParameterV2_elementNameAndIndexInExpression = `
resource "aidbox_fhir_search_parameter" "example_name" {
  name        = "first-name"
  type        = "string"
  base        = ["Patient"]
  code        = "first-name"
  expression  = "Patient.name.given"
  description = "Search patients by first name"
  url         = "https://fhir.yourcompany.com/searchparameter/first-name"
  status      = "draft"
}
`

const testAccResourceSearchParameterV2_extension = `
resource "aidbox_fhir_search_parameter" "example_extension" {
  name        = "custom-date"
  type        = "date"
  base        = ["Appointment"]
  code        = "custom-date"
  expression  = "Appointment.meta.extension.where(url = 'https://fhir.yourcompany.com/structuredefinition/custom-date').valueDateTime"
  description = "Search appointments by custom date expression"
  url         = "https://fhir.yourcompany.com/searchparameter/custom-date"
  status      = "active"
}
`

const testAccResourceSearchParameterV2_updatePhoneNumber = `
resource "aidbox_fhir_search_parameter" "example_phone" {
  name        = "phone-number"
  type        = "string"
  base        = ["Patient"]
  code        = "phone-number"
  expression  = "Patient.telecom.where(system = 'phone')"
  description = "Search patients by phone number"
  url         = "https://fhir.newdomain.com/searchparameter/phone-number"
}
`
