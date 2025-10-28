package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceGcpServiceAccount_CreateAndUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGcpServiceAccount_Create,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_gcp_service_account.test_sa", "name", "aidbox-test-sa"),
					resource.TestCheckResourceAttr("aidbox_gcp_service_account.test_sa", "service_account_email", "test-sa@my-project.iam.gserviceaccount.com"),
				),
			},
			{
				Config: testAccResourceGcpServiceAccount_Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_gcp_service_account.test_sa", "name", "aidbox-test-sa"),
					resource.TestCheckResourceAttr("aidbox_gcp_service_account.test_sa", "service_account_email", "updated-sa@my-project.iam.gserviceaccount.com"),
				),
			},
		},
	})
}

const testAccResourceGcpServiceAccount_Create = `
resource "aidbox_gcp_service_account" "test_sa" {
  name                      = "aidbox-test-sa"
  service_account_email = "test-sa@my-project.iam.gserviceaccount.com"
}
`

const testAccResourceGcpServiceAccount_Update = `
resource "aidbox_gcp_service_account" "test_sa" {
  name                      = "aidbox-test-sa"
  service_account_email = "updated-sa@my-project.iam.gserviceaccount.com"
}
`
