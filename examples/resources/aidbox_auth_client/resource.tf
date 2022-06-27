resource "aidbox_auth_client" "example" {
  secret      = "secret"
  grant_types = ["basic"]
}
