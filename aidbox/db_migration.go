package aidbox

import (
	"context"
)

// DbMigration
// Represents a migration script. Update and Delete are not supported since
// there's no trivial way to "update" or "delete" an arbitrary executed sql.
//
// Note: the /db/migrations API only handles an array, not individual resources
// with ids. Rather than exposing this API feature, we handle the wrapping/extraction
// inside an array.
//
// see https://docs.aidbox.app/modules-1/aidbox-search/usdpsql#sql-migrations
type DbMigration struct {
	Id  string `json:"id"`
	Sql string `json:"sql"`
}

func (apiClient *ApiClient) CreateDbMigration(ctx context.Context, migration *DbMigration) (*DbMigration, error) {
	response := &[]DbMigration{}
	migrations := []DbMigration{*migration}
	err := apiClient.post(ctx, migrations, "/db/migrations", response)
	if err != nil {
		return nil, err
	}

	createdMigration, err := apiClient.GetDbMigration(ctx, migration.Id)
	if err != nil {
		return nil, err
	}

	return createdMigration, nil
}

func (apiClient *ApiClient) GetDbMigration(ctx context.Context, id string) (*DbMigration, error) {
	response := &[]DbMigration{}
	err := apiClient.get(ctx, "/db/migrations", response)
	if err != nil {
		return nil, err
	}

	migration := findMigration(id, *response)
	if migration == nil {
		return nil, NotFoundError
	}

	return migration, nil
}

func findMigration(id string, migrations []DbMigration) *DbMigration {
	for _, resultMigration := range migrations {
		if resultMigration.Id == id {
			return &resultMigration
		}
	}
	return nil
}
