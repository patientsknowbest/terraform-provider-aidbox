resource "aidbox_resource" "my_aidbox_job" {
  resource = jsonencode({
    id           = "my_aidbox_job"
    resourceType = "AidboxJob"
    type         = "periodic"
  })
}