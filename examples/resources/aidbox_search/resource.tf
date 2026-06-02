resource "aidbox_search" "example_extension" {
  name         = "custom-date"
  param_parser = "reference"
  module       = "fhir-4.0.1"
  reference {
    resource_id   = "Appointment"
    resource_type = "Entity"
  }
  where = "date = '01 Jan 2014"
}
