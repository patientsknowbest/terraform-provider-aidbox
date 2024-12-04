package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceClient_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceClient_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_client.example", "id", "my-client"),
					resource.TestCheckResourceAttr("aidbox_client.example", "secret", "secret"),
					resource.TestCheckResourceAttr("aidbox_client.example", "grant_types.0", "basic"),
				),
			},
		},
	})
}

const testAccResourceClient_basic = `
resource "aidbox_client" "example" {
  name        = "my-client"
  secret      = "secret"
  grant_types = ["basic"]
}
`
