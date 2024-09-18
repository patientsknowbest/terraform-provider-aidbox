package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func resourceDbMigration() *schema.Resource {
	return &schema.Resource{
		Description: "A database migration script to be run against the db. Migrations are permanent, once created" +
			" you can't update them. You can delete the resource, but the migration will remain in the database.\n" +
			"https://docs.aidbox.app/modules-1/aidbox-search/usdpsql#sql-migrations",
		CreateContext: resourceDbMigrationCreate,
		ReadContext:   resourceDbMigrationRead,
		UpdateContext: resourceDbMigrationUpdate,
		DeleteContext: resourceDbMigrationDelete,
		Importer:      &schema.ResourceImporter{
			StateContext: resourceDbMigrationImport,
		},
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

// There's no such thing as "updating/deleting a migration". However, this doesn't
// match with terraform's resource model, so we must provide these methods.
// Trying to update will always fail with printing the below instructions to users.
// Trying to delete will always succeed, but users get a warning that deleting
// will leave some state behind in the box - as opposed to silently doing nothing
// here, which could possibly make users think they deleted the migration.
// This is OK since most of the time this deletion will occur during the removal
// of a box (i.e. not by deleting this specific resource).

func resourceDbMigrationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("Migrations cannot be updated. Add a new migration instead to achieve desired changes.")
}

func resourceDbMigrationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Warn(ctx, "**If you are deleting this box, ignore this message.**\n"+
		"Aidbox does not support deleting migrations:\n"+
		"- id '"+data.Id()+"' will be remembered and it can't be used for new migrations\n"+
		"- if you want to undo the migration script you can do this by hand")
	data.SetId("")
	return nil
}

func resourceDbMigrationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*aidbox.ApiClient)
	res, err := apiClient.GetDbMigration(ctx, d.Id(), boxIdFromData(d))
	if err != nil {
		return nil, err
	}
	mapDbMigrationToData(res, d)
	return []*schema.ResourceData{d}, nil
}

func resourceSchemaDbMigration() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// This is called name instead of id, because id is terraform-reserved, but that can't be Required
		"name": {
			Description: "Unique name for the migration, e.g. add_gin_index_to_patient",
			Type:        schema.TypeString,
			Required:    true,
		},
		"sql": {
			Description: "The sql migration script",
			Type:        schema.TypeString,
			Required:    true,
		},
	}
}
