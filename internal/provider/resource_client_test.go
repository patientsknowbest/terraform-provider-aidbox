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
					resource.TestCheckResourceAttr("aidbox_client.example", "secret", "__sha256:2BB80D537B1DA3E38BD30361AA855686BDE0EACD7162FEF6A25FE97BF527A25B"),
					resource.TestCheckResourceAttr("aidbox_client.example", "grant_types.0", "basic"),
				),
			},
		},
	})
}

const testAccResourceClient_basic = `
resource "aidbox_client" "example" {
  name        = "my-client"
  secret      = "__sha256:2BB80D537B1DA3E38BD30361AA855686BDE0EACD7162FEF6A25FE97BF527A25B"
  grant_types = ["basic"]
}
`
