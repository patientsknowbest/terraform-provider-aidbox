package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceAuthClient_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAuthClient_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_auth_client.example", "secret", "secret"),
					resource.TestCheckResourceAttr("aidbox_auth_client.example", "grant_types.0", "basic"),
				),
			},
		},
	})
}

const testAccResourceAuthClient_basic = `
resource "aidbox_auth_client" "example" {
  secret      = "secret"
  grant_types = ["basic"]
}
`
