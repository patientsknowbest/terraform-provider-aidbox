package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os/exec"
	"regexp"
	"testing"
)

// Having no real delete for the resource, CheckDestroy removes the migration
// from the db after the test. In the multibox scenario the box itself is
// deleted thus this is not required in that case.
func TestAccResourceDbMigration(t *testing.T) {
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
		CheckDestroy: func(state *terraform.State) error {
			cmd := exec.Command("sh", "./remove_test_migrations.sh")
			stdout, err := cmd.Output()
			output := string(stdout)
			if len(output) > 0 {
				t.Log(string(stdout))
			}
			return err
		},
	})
}

func TestAccResourceDbUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDbMigration,
			},
			{
				Config:      testAccResourceDbMigrationUpdate,
				ExpectError: regexp.MustCompile("Migrations cannot be updated. Add a new migration instead to achieve desired changes."),
			},
		},
		CheckDestroy: func(state *terraform.State) error {
			cmd := exec.Command("sh", "./remove_test_migrations.sh")
			stdout, err := cmd.Output()
			output := string(stdout)
			if len(output) > 0 {
				t.Log(string(stdout))
			}
			return err
		},
	})
}

func TestAccResourceDbMigration_multibox(t *testing.T) {
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

const testAccResourceDbMigrationUpdate = `
resource "aidbox_db_migration" "add_indexes" {
  name = "add_indexes"
  sql = <<-EOT
	# triggers an update which will fail
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
