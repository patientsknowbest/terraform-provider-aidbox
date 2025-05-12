package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceIdentityProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIdentityProvider,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "title", "MyIDP"),
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "userinfo_source", "id-token"),
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "client.0.id", "some_client_id"),
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "client.0.secret", "some_client_secret"),
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "scopes.0", "https://www.myidp.com/scope1"),
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "scopes.1", "https://www.myidp.com/scope2"),
				),
			},
			{
				Config: testAccResourceIdentityProvider_updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "title", "MyIDP"),
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "userinfo_source", "userinfo-endpoint"),
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "client.0.id", "some_client_id"),
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "client.0.secret", "some_client_secret_updated"),
					resource.TestCheckResourceAttr("aidbox_identity_provider.myidp", "scopes.0", "https://www.myidp.com/scope1"),
					resource.TestCheckNoResourceAttr("aidbox_identity_provider.myidp", "scopes.1"),
				),
			},
		},
	})
}

const testAccResourceIdentityProvider = `
resource "aidbox_identity_provider" "myidp" {
  title = "MyIDP"
  userinfo_source = "id-token"

  scopes = [
    "https://www.myidp.com/scope1",
    "https://www.myidp.com/scope2",
  ]

  client {
    id = "some_client_id"
    secret = "some_client_secret"
  }
}
`

const testAccResourceIdentityProvider_updated = `
resource "aidbox_identity_provider" "myidp" {
  title = "MyIDP"
  userinfo_source = "userinfo-endpoint"

  scopes = [
    "https://www.myidp.com/scope1",
  ]

  client {
    id = "some_client_id"
    secret = "some_client_secret_updated"
  }
}
`
