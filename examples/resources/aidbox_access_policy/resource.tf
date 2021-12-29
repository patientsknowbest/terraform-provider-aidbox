resource "aidbox_access_policy" "example" {
  description = "A policy to allow postman to access data"
  engine      = "json-schema"
  schema = jsonencode({
    "required" = [
      "client",
      "uri",
    "request-method"]
    "properties" = {
      "uri" = {
        "type"    = "string"
        "pattern" = "^/fhir/.*"
      }
      "client" = {
        "required" = ["id"]
        "properties" = {
          "id" = {
            const = "postman"
          }
        }
      }
      "request-method" = {
        "const" = "get"
      }
    }
  })
}