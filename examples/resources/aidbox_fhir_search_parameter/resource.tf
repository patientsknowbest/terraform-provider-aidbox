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