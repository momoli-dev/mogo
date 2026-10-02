// Package main is an example of using mogosql.
package main

import (
	"context"
	"embed"

	"github.com/momoli-dev/mogo/mogosql"
)

//go:embed migrations/*.sql
var MigrationFS embed.FS

func main() {
	ctx := context.Background()

	// Migrate a Postgres database with embedded migrations.
	//nolint:gosec // Hardcoded credential for the sake of example.
	if err := mogosql.MigratePostgres(ctx, mogosql.PostgresMigrateParams{
		ConnURL:      "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		MigrationsFS: MigrationFS,
	}); err != nil {
		panic(err)
	}
}
