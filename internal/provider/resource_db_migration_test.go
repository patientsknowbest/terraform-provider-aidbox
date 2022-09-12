package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// TODO remove .Skip calls from tests
// - remove state from the db after test run
// - ignore error thrown by post-script due to trying to delete resources
// see https://pkbdev.atlassian.net/browse/PHR-11231

func TestAccResourceDbMigration(t *testing.T) {
	t.Skip("This test is intended to be run only by hand. It's skipped because migrations can't be deleted and the post script for the test would fail, leaving unwanted state in the db.")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDbMigration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_db_migration.add_indexes", "id", "add_indexes"),
					resource.TestCheckResourceAttr("aidbox_db_migration.add_indexes", "name", "add_indexes"),
					resource.TestCheckResourceAttr("aidbox_db_migration.add_indexes", "sql", "CREATE INDEX appointment_resource_idx ON public.appointment USING gin (resource);\nCREATE INDEX patient_resource_idx ON public.patient USING gin (resource);\n"),
				),
			},
		},
	})
}

func TestAccResourceDbMigration_multibox(t *testing.T) {
	t.Skip("This test is intended to be run only by hand. It's skipped because migrations can't be deleted and the post script for the test would fail, leaving unwanted state in the db.")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testMultiboxProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDbMigrationMultibox,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_db_migration.add_indexes_on_multibox", "id", "add_indexes_on_multibox"),
					resource.TestCheckResourceAttr("aidbox_db_migration.add_indexes_on_multibox", "name", "add_indexes_on_multibox"),
					resource.TestCheckResourceAttr("aidbox_db_migration.add_indexes_on_multibox", "sql", "CREATE INDEX patient_txid_idx ON patient(txid);\nCREATE INDEX person_txid_idx ON person(txid);\n"),
				),
			},
		},
	})
}

const testAccResourceDbMigration = `
resource "aidbox_db_migration" "add_indexes" {
  name = "add_indexes"
  sql = <<-EOT
	CREATE INDEX appointment_resource_idx ON public.appointment USING gin (resource);
	CREATE INDEX patient_resource_idx ON public.patient USING gin (resource);
  EOT
}
`

const testAccResourceDbMigrationMultibox = `
resource "aidbox_box" "mybox" {
  name = "mybox"
  fhir_version  = "fhir-3.0.1"
  description = "A box instance within multibox, a multi-tenant aidbox server"
}
resource "aidbox_db_migration" "add_indexes_on_multibox" {
  box_id = aidbox_box.mybox.name
  name = "add_indexes_on_multibox"
  sql = <<-EOT
	CREATE INDEX patient_txid_idx ON patient(txid);
	CREATE INDEX person_txid_idx ON person(txid);
  EOT
}
`
