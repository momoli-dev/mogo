// Package mogosqltest provides ready to use, parallelizable SQL databases for tests.
package mogosqltest

import (
	"context"
	"io/fs"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/goosemigrator"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgres creates a migrated test database and registers connection cleanup.
func NewPostgres(t *testing.T, migrationsFS fs.FS) *pgx.Conn {
	t.Helper()

	conn, _ := NewPostgresWithConnURL(t, migrationsFS)

	return conn
}

// NewPostgresWithConnURL also returns the test database connection URL.
func NewPostgresWithConnURL(t *testing.T, migrationsFS fs.FS) (*pgx.Conn, string) {
	t.Helper()
	ctx := t.Context()

	const (
		postgresUser         = "postgres"
		postgresPassword     = "postgres"
		postgresHost         = "localhost"
		postgresPort         = "5432"
		postgresCloseTimeout = 5 * time.Second
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

	conn, err := pgx.Connect(ctx, conf.URL())
	require.NoError(t, err)
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.WithoutCancel(t.Context()), postgresCloseTimeout)
		defer cancel()

		err := conn.Close(closeCtx)
		require.NoError(t, err)
	})

	return conn, conf.URL()
}
