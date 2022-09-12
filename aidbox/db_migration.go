package aidbox

import (
	"context"
	"errors"
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

func (apiClient *ApiClient) CreateDbMigration(ctx context.Context, migration *DbMigration, boxId string) (*DbMigration, error) {
	response := &[]DbMigration{}
	migrations := []DbMigration{*migration}
	err := apiClient.post(ctx, migrations, "/db/migrations", boxId, response)
	if err != nil {
		return nil, err
	}

	// Response might return other ongoing migrations, hence the search
	createdMigration := findMigration(migration.Id, *response)
	if createdMigration == nil {
		return nil, errors.New("failed to create migration with id: " + migration.Id + ", response did not contain the requested migration id")
	}

	return createdMigration, nil
}

func (apiClient *ApiClient) GetDbMigration(ctx context.Context, id, boxId string) (*DbMigration, error) {
	response := &[]DbMigration{}
	err := apiClient.get(ctx, "/db/migrations", boxId, response)
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
