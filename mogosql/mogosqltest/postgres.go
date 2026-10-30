// Package mogosqltest provides ready to use, parallelizable SQL databases for tests.
package mogosqltest

import (
	"io/fs"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/momoli-dev/mogo/mogosql"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/goosemigrator"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgres returns an isolated, migrated pool and its URL, and registers cleanup.
func NewPostgres(t *testing.T, migrationsFS fs.FS) (*pgxpool.Pool, string) {
	t.Helper()
	ctx := t.Context()

	const (
		postgresUser     = "postgres"
		postgresPassword = "postgres"
		postgresHost     = "localhost"
		postgresPort     = "5432"
	)

	dbconf := pgtestdb.Config{
		DriverName:                "pgx",
		User:                      postgresUser,
		Password:                  postgresPassword,
		Host:                      postgresHost,
		Port:                      postgresPort,
		Options:                   "sslmode=disable",
		Database:                  "postgres",
		TestRole:                  nil,
		ForceTerminateConnections: false,
	}

	m := goosemigrator.New(
		"migrations",
		goosemigrator.WithFS(migrationsFS),
	)

	conf := pgtestdb.Custom(t, dbconf, m)
	require.NotEqual(t, dbconf, conf)

	pool, err := mogosql.NewPostgresPool(ctx, mogosql.PostgresPoolParams{
		ConnStr:    conf.URL(),
		HasPostgis: false,
	})
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	return pool, conf.URL()
}
