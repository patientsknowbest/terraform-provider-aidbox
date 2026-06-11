package provider

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAidboxResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAidboxResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_resource.my_resource", "id", "AidboxJob/my_aidbox_job"),
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "resource", compIgnoreMeta(my_resource)),
				),
			},
			{
				Config: testAccResourceAidboxResource_updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_resource.my_resource", "id", "AidboxJob/my_aidbox_job"),
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "resource", compIgnoreMeta(my_resource_updated)),
				),
			},
		},
	})
}

func compIgnoreMeta(expected string) func(string) error {
	return func(actual string) error {
		if jsonDiffSuppressMetaFunc("", expected, actual, nil) {
			return nil
		} else {
			return errors.New("resource [" + actual + "] did not match expected [" + expected + "]")
		}
	}
}

const testAccResourceAidboxResource = `
resource "aidbox_resource" "my_resource" {
  resource = <<-EOT` + my_resource + `
EOT
}

`
const my_resource = `
{
  "id": "my_aidbox_job",
  "resourceType": "AidboxJob"
}
`

const testAccResourceAidboxResource_updated = `
resource "aidbox_resource" "my_resource" {
  resource = <<-EOT` + my_resource_updated + `
EOT
}
`

const my_resource_updated = `
{
  "id": "my_aidbox_job",
  "resourceType": "AidboxJob",
  "type": "periodic"
}
`
