resource "aidbox_gcp_service_account" "default_gcp_account" {
  name                  = "aidbox-rc"
  service_account_email = "sa-email@my-project.iam.gserviceaccount.com"
}