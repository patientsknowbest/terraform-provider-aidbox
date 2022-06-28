resource "aidbox_auth_client" "example" {
  name        = "my-client"
  secret      = "secret"
  grant_types = ["basic"]
}
