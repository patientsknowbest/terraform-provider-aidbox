package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceSearchParameter_elementNameAndPatternFilterInExpression(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearchParameter_elementNameAndPatternFilterInExpression,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "id", "Patient.phone-number"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "name", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "type", "string"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "base.0", "Patient"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "code", "phone-number"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "expression", "Patient.telecom.where(system = 'phone')"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "description", "Search patients by phone number"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "url", "https://fhir.yourcompany.com/searchparameter/phone-number"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_phone", "status", "active"),
				),
			},
		},
	})
}

func TestAccResourceSearchParameter_elementNameAndIndexInExpression(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearchParameter_elementNameAndIndexInExpression,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "id", "Patient.first-name"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "name", "first-name"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "type", "string"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "base.0", "Patient"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "code", "first-name"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "expression", "Patient.name.given"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "description", "Search patients by first name"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "url", "https://fhir.yourcompany.com/searchparameter/first-name"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_name", "status", "draft"),
				),
			},
		},
	})
}

func TestAccResourceSearchParameter_extension(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSearchParameter_extension,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "id", "Appointment.custom-date"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "name", "custom-date"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "type", "date"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "base.0", "Appointment"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "code", "custom-date"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "expression", "Appointment.meta.extension.where(url = 'https://fhir.yourcompany.com/structuredefinition/custom-date').valueDateTime"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "description", "Search appointments by custom date expression"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "url", "https://fhir.yourcompany.com/searchparameter/custom-date"),
					resource.TestCheckResourceAttr("aidbox_search_parameter.example_extension", "status", "active"),
				),
			},
		},
	})
}

const testAccResourceSearchParameter_elementNameAndPatternFilterInExpression = `
resource "aidbox_search_parameter" "example_phone" {
  name        = "phone-number"
  type        = "string"
  base        = ["Patient"]
  code        = "phone-number"
  expression  = "Patient.telecom.where(system = 'phone')"
  description = "Search patients by phone number"
  url         = "https://fhir.yourcompany.com/searchparameter/phone-number"
}
`

const testAccResourceSearchParameter_elementNameAndIndexInExpression = `
resource "aidbox_search_parameter" "example_name" {
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

const testAccResourceSearchParameter_extension = `
resource "aidbox_search_parameter" "example_extension" {
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
