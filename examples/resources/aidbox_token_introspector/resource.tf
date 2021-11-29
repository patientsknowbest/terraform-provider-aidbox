resource "aidbox_token_introspector" "example" {
  jwks_uri = "http://keycloak:8080/auth/realms/pkb/protocol/openid-connect/certs"
  jwt {
    iss = "http://keycloak:8080/auth/realms/pkb"
  }
}