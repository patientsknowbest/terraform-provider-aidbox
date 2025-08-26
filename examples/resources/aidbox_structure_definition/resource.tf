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
  differential    = <<-EOT
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
EOT
}