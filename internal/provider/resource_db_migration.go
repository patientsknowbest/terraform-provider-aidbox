package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceDbMigration() *schema.Resource {
	return &schema.Resource{
		Description: "A database migration script to be run against the db. Migrations are permanent, once created" +
			" you can't update/delete them. https://docs.aidbox.app/modules-1/aidbox-search/usdpsql#sql-migrations",
		CreateContext: resourceDbMigrationCreate,
		ReadContext:   resourceDbMigrationRead,
		DeleteContext: resourceDbMigrationDelete,
		Schema:        resourceFullSchema(resourceSchemaDbMigration()),
	}
}

func mapDbMigrationToData(migration *aidbox.DbMigration, data *schema.ResourceData) {
	data.SetId(migration.Id)
	data.Set("name", migration.Id)
	data.Set("sql", migration.Sql)
}

func mapDbMigrationFromData(data *schema.ResourceData) *aidbox.DbMigration {
	return &aidbox.DbMigration{
		Id:  data.Get("name").(string),
		Sql: data.Get("sql").(string),
	}
}

func resourceDbMigrationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	migration := mapDbMigrationFromData(data)
	result, err := apiClient.CreateDbMigration(ctx, migration, boxIdFromData(data))
	if err != nil {
		return diag.FromErr(err)
	}
	mapDbMigrationToData(result, data)
	return nil
}

func resourceDbMigrationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*aidbox.ApiClient)
	result, err := apiClient.GetDbMigration(ctx, data.Id(), boxIdFromData(data))
	if err != nil {
		if handleNotFoundError(err, data) {
			return nil
		}
		return diag.FromErr(err)
	}
	mapDbMigrationToData(result, data)
	return nil
}

// There's no such thing as "deleting a migration". However, this doesn't match
// with terraform's resource model, so we must provide a delete method.
// Trying to delete will always fail: this warns users that what they're trying
// to do is impossible - as opposed to silently doing nothing here, which could
// possibly make users think they deleted the migration.
func resourceDbMigrationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("Migrations cannot be deleted. To delete a migration, undo it from the database manually," +
		" then remove the state from terraform. NB if not entirely cleared down, aidbox will remember the id/script")
}

func resourceSchemaDbMigration() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// This is called name instead of id, because id is terraform-reserved, but that can't be Required
		"name": {
			ForceNew:    true,
			Description: "Unique name for the migration, e.g. add_gin_index_to_patient",
			Type:        schema.TypeString,
			Required:    true,
		},
		"sql": {
			ForceNew:    true,
			Description: "The sql migration script",
			Type:        schema.TypeString,
			Required:    true,
		},
	}
}
