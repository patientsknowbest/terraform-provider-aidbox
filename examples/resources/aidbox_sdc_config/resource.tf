resource "aidbox_sdc_config" "default_storage_config" {
  name    = "forms-storage"
  default = true

  storage {
    bucket = "attachment-store-rc"
    account {
      reference = "GcpServiceAccount/aidbox-rc"
    }
  }
}