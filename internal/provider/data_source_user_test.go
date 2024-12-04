package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUser("data.aidbox_user.admin"),
				),
			},
		},
	})
}

func testAccCheckUser(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		twoFactorEnabled := rs.Primary.Attributes["two_factor_enabled"]

		if twoFactorEnabled != "false" {
			return fmt.Errorf("expected user with ID %s to have two_factor_enabled=false but got %s ", id, twoFactorEnabled)
		}

		return nil
	}
}

const testAccDataSourceUser = `
data "aidbox_user" "admin" {
  id = "admin"
}
`
