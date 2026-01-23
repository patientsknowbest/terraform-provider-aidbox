package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccResourceStructureDefinitionOverride_setupDefaultPatientProfileThenUpdate(t *testing.T) {
	ptProfV1, err := os.ReadFile("./test_resources/patient-profile-nhs.json")
	if err != nil {
		t.Fatal(err)
	}
	ptProfV1String := string(ptProfV1)
	ptProfV2, err := os.ReadFile("./test_resources/patient-profile-nhs-mandatory-identifier.json")
	if err != nil {
		t.Fatal(err)
	}
	ptProfV2String := string(ptProfV2)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStructureDefinitionOverride_setupDefaultPatientProfileThenUpdate_setup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_structure_definition_override.patient_profile", "url", "http://hl7.org/fhir/StructureDefinition/Patient"),
					resource.TestCheckResourceAttrSet("aidbox_structure_definition_override.patient_profile", "structure_definition_override"),
					resource.TestCheckResourceAttrWith("aidbox_structure_definition_override.patient_profile", "structure_definition_override", func(valueFromServer string) error {
						assert.True(t, jsonDiffSuppressFunc("", ptProfV1String, valueFromServer, nil), "Value received from server does not match (semantically): %s", valueFromServer)
						return nil
					}),
				),
			},
			{
				Config: testAccResourceStructureDefinitionOverride_setupDefaultPatientProfileThenUpdate_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_structure_definition_override.patient_profile", "url", "http://hl7.org/fhir/StructureDefinition/Patient"),
					resource.TestCheckResourceAttrSet("aidbox_structure_definition_override.patient_profile", "structure_definition_override"),
					resource.TestCheckResourceAttrWith("aidbox_structure_definition_override.patient_profile", "structure_definition_override", func(valueFromServer string) error {
						assert.True(t, jsonDiffSuppressFunc("", ptProfV2String, valueFromServer, nil), "Value received from server does not match (semantically): %s", valueFromServer)
						return nil
					}),
				),
			},
		},
	})
}

const testAccResourceStructureDefinitionOverride_setupDefaultPatientProfileThenUpdate_setup = `
resource "aidbox_structure_definition_override" "patient_profile" {
  url = "http://hl7.org/fhir/StructureDefinition/Patient"
  structure_definition_override = file("${path.module}/test_resources/patient-profile-nhs.json")
}
`

const testAccResourceStructureDefinitionOverride_setupDefaultPatientProfileThenUpdate_update = `
resource "aidbox_structure_definition_override" "patient_profile" {
  url = "http://hl7.org/fhir/StructureDefinition/Patient"
  structure_definition_override = file("${path.module}/test_resources/patient-profile-nhs-mandatory-identifier.json")
}
`
