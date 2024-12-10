resource "aidbox_db_migration" "add_gin_index_on_appointment_and_patient" {
  name = "add_gin_index_on_appointment_and_patient"
  sql  = <<-EOT
    CREATE INDEX appointment_resource_idx ON public.appointment USING gin (resource);
    CREATE INDEX patient_resource_idx ON public.patient USING gin (resource);
  EOT
}