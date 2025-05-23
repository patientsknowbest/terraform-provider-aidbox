package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os/exec"
	"regexp"
	"testing"
)

// Having no real delete for the resource, CheckDestroy removes the migration
// from the db after the test.
func TestAccResourceDbMigration_Create(t *testing.T) {
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

func TestAccResourceDbMigration_MultipleCreates(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDbMigration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_db_migration.add_indexes", "id", "add_indexes"),
				),
			},
			{
				Config: testAccResourceDbMigration2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_db_migration.add_index_on_practitioner", "id", "add_index_on_practitioner"),
				),
			},
			{
				Config: testAccResourceDbMigration3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_db_migration.add_test_table", "id", "add_test_table"),
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

func TestAccResourceDbMigration_UpdateAlwaysErrors(t *testing.T) {
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

const testAccResourceDbMigration = `
resource "aidbox_db_migration" "add_indexes" {
  name = "add_indexes"
  sql = <<-EOT
	CREATE INDEX appointment_resource_idx ON public.appointment USING gin (resource);
	CREATE INDEX patient_resource_idx ON public.patient USING gin (resource);
  EOT
}
`

const testAccResourceDbMigration2 = `
resource "aidbox_db_migration" "add_index_on_practitioner" {
  name = "add_index_on_practitioner"
  sql = <<-EOT
    CREATE INDEX practitioner_txid_idx on practitioner (txid);
  EOT
}
`

const testAccResourceDbMigration3 = `
resource "aidbox_db_migration" "add_test_table" {
  name = "add_test_table"
  sql = <<-EOT
	CREATE TABLE migration_test(
	   id serial PRIMARY KEY,
	   whatever VARCHAR (255) UNIQUE NOT NULL
	);
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
