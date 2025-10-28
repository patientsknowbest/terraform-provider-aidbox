package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceSDCConfig_CreateAndUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { requireSchemaMode(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSDCConfig_Create,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "name", "forms-storage"),
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "default", "true"),
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "description", "Default SDC config for forms"),
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "storage.0.bucket", "attachment-store-test"),
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "storage.0.account.0.reference", "GcpServiceAccount/aidbox-test"),
				),
			},
			{
				Config: testAccResourceSDCConfig_Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "name", "forms-storage"),
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "default", "false"),
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "description", "Updated SDC config"),
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "storage.0.bucket", "attachment-store-test"),
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "storage.0.account.0.reference", "GcpServiceAccount/aidbox-test"),
				),
			},
		},
	})
}

const testAccResourceSDCConfig_Create = `
resource "aidbox_gcp_service_account" "test_account" {
  name                      = "aidbox-test"
  service_account_email = "test-for-sdc@my-project.iam.gserviceaccount.com"
}

resource "aidbox_sdc_config" "default_storage" {
  name        = "forms-storage"
  default     = true
  description = "Default SDC config for forms"

  storage {
    bucket = "attachment-store-test"
    account {
      reference = "GcpServiceAccount/${aidbox_gcp_service_account.test_account.id}"
    }
  }
}
`

const testAccResourceSDCConfig_Update = `
resource "aidbox_gcp_service_account" "test_account" {
  name                      = "aidbox-test"
  service_account_email = "test-for-sdc@my-project.iam.gserviceaccount.com"
}

resource "aidbox_sdc_config" "default_storage" {
  name        = "forms-storage"
  default     = false 
  description = "Updated SDC config" 

  storage {
    bucket = "attachment-store-test"
    account {
      reference = "GcpServiceAccount/${aidbox_gcp_service_account.test_account.id}"
    }
  }
}
`
