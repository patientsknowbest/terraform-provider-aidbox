package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccResourceStructureDefinition_setupAndUpdatePatientProfile(t *testing.T) {
	previousIdState := ""
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStructureDefinition_setupPatientProfile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "name", "patient-profile"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "url", "https://fhir.yourcompany.com/structuredefinition/patient"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "base_definition", "http://hl7.org/fhir/StructureDefinition/Patient"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "derivation", "constraint"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "abstract", "false"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "type", "Patient"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "status", "active"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "kind", "resource"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "version", "0.0.1"),
					resource.TestCheckResourceAttrWith("aidbox_structure_definition.patient_profile", "differential", func(valueFromServer string) error {
						assert.True(t, jsonDiffSuppressFunc("", differential, valueFromServer, nil), "Value received from server does not match (semantically): %s", valueFromServer)
						return nil
					}),
					resource.TestCheckResourceAttrWith("aidbox_structure_definition.patient_profile", "id", func(id string) error {
						previousIdState = id
						return nil
					}),
				),
			},
			{
				Config: testAccResourceStructureDefinition_updatePatientProfile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("aidbox_structure_definition.patient_profile", "id", &previousIdState),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "name", "patient-profile"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "url", "https://fhir.yourcompany.com/structuredefinition/patient"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "base_definition", "http://hl7.org/fhir/StructureDefinition/Patient"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "derivation", "constraint"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "abstract", "false"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "type", "Patient"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "status", "active"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "kind", "resource"),
					resource.TestCheckResourceAttr("aidbox_structure_definition.patient_profile", "version", "0.0.1"),
					resource.TestCheckResourceAttrWith("aidbox_structure_definition.patient_profile", "differential", func(valueFromServer string) error {
						assert.True(t, jsonDiffSuppressFunc("", differential_updated, valueFromServer, nil), "Value received from server does not match (semantically): %s", valueFromServer)
						return nil
					}),
				),
			},
		},
	})
}

const differential = `
    {
      "element": [
        {
          "id": "Patient",
          "path": "Patient",
          "constraint": [
            {
              "key": "unique-system",
              "severity": "error",
              "human": "System must be unique among identifiers",
              "expression": "Patient.identifier.system.count() = Patient.identifier.system.distinct().count()"
            }
          ]
        },
        {
          "id": "Patient.identifier",
          "path": "Patient.identifier",
          "min": 1
        },
        {
          "id": "Patient.identifier.system",
          "path": "Patient.identifier.system",
          "min": 1
        },
        {
          "id": "Patient.identifier.value",
          "path": "Patient.identifier.value",
          "min": 1
        },
        {
          "id": "Patient.managingOrganization",
          "path": "Patient.managingOrganization",
          "min": 1
        },
        {
          "id": "Patient.managingOrganization.reference",
          "path": "Patient.managingOrganization.reference",
          "min": 1
        }
      ]
    }
`

const testAccResourceStructureDefinition_setupPatientProfile = `
resource "aidbox_structure_definition" "patient_profile" {
  name            = "patient-profile"
  url             = "https://fhir.yourcompany.com/structuredefinition/patient"
  base_definition = "http://hl7.org/fhir/StructureDefinition/Patient"
  derivation      = "constraint"
  abstract        = false
  type            = "Patient"
  status          = "active"
  kind            = "resource"
  version         = "0.0.1"
  differential    = <<-EOT` + differential + `
EOT
}
`

const differential_updated = `
    {
      "element": [
        {
          "id": "Patient",
          "path": "Patient",
          "constraint": [
            {
              "key": "unique-system",
              "severity": "error",
              "human": "System must be unique among identifiers",
              "expression": "Patient.identifier.system.count() = Patient.identifier.system.distinct().count()"
            }
          ]
        }
      ]
    }
`

const testAccResourceStructureDefinition_updatePatientProfile = `
resource "aidbox_structure_definition" "patient_profile" {
  name            = "patient-profile"
  url             = "https://fhir.yourcompany.com/structuredefinition/patient"
  base_definition = "http://hl7.org/fhir/StructureDefinition/Patient"
  derivation      = "constraint"
  abstract        = false
  type            = "Patient"
  status          = "active"
  kind            = "resource"
  version         = "0.0.1"
  differential    = <<-EOT` + differential_updated + `
EOT
}
`
