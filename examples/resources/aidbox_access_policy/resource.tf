resource "aidbox_access_policy" "example_schema" {
  description = "A policy to allow postman to access data"
  engine      = "json-schema"
  schema = jsonencode({
    "required" = [
      "client",
      "uri",
      "request-method"
    ]
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

resource "aidbox_access_policy" "example_allow" {
  description = "A policy to allow client to access data"
  engine      = "allow"
  link {
    resource_id   = "client-id"
    resource_type = "Client"
  }
}
