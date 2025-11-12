package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestAccResourceGcpServiceAccount_CreateAndUpdate(t *testing.T) {
	previousIdState := ""
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGcpServiceAccount_Create,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_gcp_service_account.test_sa", "name", "aidbox-test-sa"),
					resource.TestCheckResourceAttr("aidbox_gcp_service_account.test_sa", "service_account_email", "test-sa@my-project.iam.gserviceaccount.com"),
					resource.TestCheckResourceAttr("aidbox_gcp_service_account.test_sa", "private_key", "fake-private-key-v1"),
					resource.TestCheckResourceAttrWith("aidbox_gcp_service_account.test_sa", "id", func(id string) error {
						previousIdState = id
						return nil
					}),
				),
			},
			{
				Config: testAccResourceGcpServiceAccount_Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_gcp_service_account.test_sa", "name", "aidbox-test-sa"),
					resource.TestCheckResourceAttr("aidbox_gcp_service_account.test_sa", "service_account_email", "updated-sa@my-project.iam.gserviceaccount.com"),
					resource.TestCheckResourceAttr("aidbox_gcp_service_account.test_sa", "private_key", "fake-private-key-v2"),
					resource.TestCheckResourceAttrWith("aidbox_gcp_service_account.test_sa", "id", func(id string) error {
						assert.Equalf(t, previousIdState, id, "Resource logical id unexpectedly changed after resource update")
						return nil
					}),
				),
			},
		},
	})
}

const testAccResourceGcpServiceAccount_Create = `
resource "aidbox_gcp_service_account" "test_sa" {
  name                      = "aidbox-test-sa"
  service_account_email = "test-sa@my-project.iam.gserviceaccount.com"
  private_key             = "fake-private-key-v1"
}
`

const testAccResourceGcpServiceAccount_Update = `
resource "aidbox_gcp_service_account" "test_sa" {
  name                      = "aidbox-test-sa"
  service_account_email = "updated-sa@my-project.iam.gserviceaccount.com"
  private_key             = "fake-private-key-v2"
}
`
