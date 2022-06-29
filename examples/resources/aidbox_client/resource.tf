resource "aidbox_client" "example" {
  name        = "my-client"
  secret      = "secret"
  grant_types = ["basic"]
}
