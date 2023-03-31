resource "aidbox_search_parameter" "example_extension" {
  name   = "custom-date"
  module = "fhir-4.0.1"
  type   = "date"
  reference {
    resource_id   = "Appointment"
    resource_type = "Entity"
  }
  expression = [
    "meta", "extension", "url|https://fhir.patientsknowbest.com/structuredefinition/custom-date", "valueDateTime"
  ]
}