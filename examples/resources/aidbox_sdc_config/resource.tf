resource "aidbox_sdcconfig" "default_storage_config" {
  name    = "forms-storage"
  default = true

  storage {
    bucket = "attachment-store-rc"
    account {
      reference = "GcpServiceAccount/aidbox-rc"
    }
  }
}