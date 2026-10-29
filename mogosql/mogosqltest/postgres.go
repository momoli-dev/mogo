// Package mogosqltest provides ready to use, parallelizable SQL databases for tests.
package mogosqltest

import (
	"io/fs"
	"testing"

	"github.com/momoli-dev/mogo/mogosql"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/goosemigrator"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgres creates an isolated, migrated test database and registers pool cleanup.
func NewPostgres(t *testing.T, migrationsFS fs.FS) *mogosql.PostgresDB {
	t.Helper()

	conn, _ := NewPostgresWithConnURL(t, migrationsFS)

	return conn
}

// NewPostgresWithConnURL also returns the test database connection URL.
func NewPostgresWithConnURL(t *testing.T, migrationsFS fs.FS) (*mogosql.PostgresDB, string) {
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

	conn, err := mogosql.NewPostgresDB(ctx, mogosql.PostgresDBParams{
		ConnStr:    conf.URL(),
		HasPostgis: false,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		err := conn.Close()
		require.NoError(t, err)
	})

	return conn, conf.URL()
}
