resource "aidbox_box" "my_box" {
  id           = "my-box"
  fhir_version = "fhir-3.0.1"
  description  = "A box instance within multibox, a multi-tenant aidbox server"
  env = [
    "ENV1=foo",
    "ENV2=bar"
  ]
}