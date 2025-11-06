package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
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
					resource.TestCheckResourceAttrWith("aidbox_sdc_config.default_storage", "storage", func(valueFromServer string) error {
						assert.True(t, jsonDiffSuppressFunc("", storage_v1, valueFromServer, nil), "Value received from server does not match (semantically): %s", valueFromServer)
						return nil
					}),
				),
			},
			{
				Config: testAccResourceSDCConfig_Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "name", "forms-storage"),
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "default", "false"),
					resource.TestCheckResourceAttr("aidbox_sdc_config.default_storage", "description", "Updated SDC config"),
					resource.TestCheckResourceAttrWith("aidbox_sdc_config.default_storage", "storage", func(valueFromServer string) error {
						assert.True(t, jsonDiffSuppressFunc("", storage_v2, valueFromServer, nil), "Value received from server does not match (semantically): %s", valueFromServer)
						return nil
					}),
				),
			},
		},
	})
}

const storage_v1 = `
{
  "bucket": "attachment-store-rc",
  "account": {
    "reference": "GcpServiceAccount/aidbox-test-for-sdc"
  }
}
`

const testAccResourceSDCConfig_Create = `
resource "aidbox_gcp_service_account" "test_account" {
  name       = "aidbox-test-for-sdc"
  service_account_email = "test-sa-email@example.com"
}

resource "aidbox_sdc_config" "default_storage" {
  name        = "forms-storage"
  default     = true
  description = "Default SDC config for forms"
  storage     = <<-EOT
  {
    "bucket": "attachment-store-rc",
    "account": {
      "reference": "GcpServiceAccount/${aidbox_gcp_service_account.test_account.id}"
    }
  }
  EOT
}
`

const storage_v2 = `
{
  "bucket": "attachment-store-rc-UPDATED",
  "account": {
    "reference": "GcpServiceAccount/aidbox-test-for-sdc"
  }
}
`

const testAccResourceSDCConfig_Update = `
resource "aidbox_gcp_service_account" "test_account" {
  name       = "aidbox-test-for-sdc"
  service_account_email = "test-sa-email@example.com"
}

resource "aidbox_sdc_config" "default_storage" {
  name        = "forms-storage"
  default     = false 
  description = "Updated SDC config" 
  storage     = <<-EOT
  {
    "bucket": "attachment-store-rc-UPDATED",
    "account": {
      "reference": "GcpServiceAccount/${aidbox_gcp_service_account.test_account.id}"
    }
  }
  EOT
}
`
