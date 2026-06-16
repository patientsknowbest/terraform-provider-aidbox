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
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "resource", compIgnoreJsonDiff(my_resource)),
				),
			},
			{
				Config: testAccResourceAidboxResource_updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_resource.my_resource", "id", "AidboxJob/my_aidbox_job"),
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "resource", compIgnoreJsonDiff(my_resource_updated)),
				),
			},
			{
				Config: testAccResourceAidboxResource_updated_newid,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_resource.my_resource", "id", "AidboxJob/my_aidbox_job_newid"),
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "resource", compIgnoreJsonDiff(my_resource_updated_newid)),
				),
			},
		},
	})
}

func compIgnoreJsonDiff(expected string) func(string) error {
	return func(actual string) error {
		if jsonDiffSuppressFunc("", expected, actual, nil) {
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

const testAccResourceAidboxResource_updated_newid = `
resource "aidbox_resource" "my_resource" {
  resource = <<-EOT` + my_resource_updated_newid + `
EOT
}
`

const my_resource_updated_newid = `
{
  "id": "my_aidbox_job_newid",
  "resourceType": "AidboxJob",
  "type": "periodic"
}
`

func TestAccResourceAidboxResourceNoId(t *testing.T) {
	var assignedId string
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAidboxResource_noid,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "id", func(value string) error {
						assignedId = value
						return nil
					}),
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "resource", compIgnoreJsonDiff(my_resource_noid)),
				),
			},
			{
				Config: testAccResourceAidboxResource_noid_updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("aidbox_resource.my_resource", "id", &assignedId),
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "resource", compIgnoreJsonDiff(my_resource_noid_updated)),
				),
			},
			{
				Config: testAccResourceAidboxResource_noid_updated_changed_resourcetype,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "id", func(value string) error {
						if value == assignedId {
							return errors.New("new ID should have been assigned")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "resource", compIgnoreJsonDiff(my_resource_noid_updated_resourcetype)),
				),
			},
			{
				Config: testAccResourceAidboxResource_noid_updated_addid,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_resource.my_resource", "id", "AidboxConfig/my_config_id"),
					resource.TestCheckResourceAttrWith("aidbox_resource.my_resource", "resource", compIgnoreJsonDiff(my_resource_noid_updated_addid)),
				),
			},
		},
	})
}

const testAccResourceAidboxResource_noid = `
resource "aidbox_resource" "my_resource" {
  resource = <<-EOT` + my_resource_noid + `
EOT
}
`
const my_resource_noid = `
{
  "resourceType": "AidboxJob"
}
`

const testAccResourceAidboxResource_noid_updated = `
resource "aidbox_resource" "my_resource" {
  resource = <<-EOT` + my_resource_noid_updated + `
EOT
}
`
const my_resource_noid_updated = `
{
  "resourceType": "AidboxJob",
  "type": "periodic"
}
`
const testAccResourceAidboxResource_noid_updated_changed_resourcetype = `
resource "aidbox_resource" "my_resource" {
  resource = <<-EOT` + my_resource_noid_updated_resourcetype + `
EOT
}
`
const my_resource_noid_updated_resourcetype = `
{
  "resourceType": "AidboxConfig"
}
`

const testAccResourceAidboxResource_noid_updated_addid = `
resource "aidbox_resource" "my_resource" {
  resource = <<-EOT` + my_resource_noid_updated_addid + `
EOT
}
`
const my_resource_noid_updated_addid = `
{
  "id": "my_config_id",
  "resourceType": "AidboxConfig"
}
`
