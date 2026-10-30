package main_test

import (
	"testing"

	example "github.com/momoli-dev/mogo/example/mogosql"
	"github.com/momoli-dev/mogo/mogosql/mogosqltest"
	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {
	t.Parallel()
	pool, _ := mogosqltest.NewPostgres(t, example.MigrationFS)
	var name string
	require.NoError(t, pool.QueryRow(t.Context(), "SELECT current_database()").Scan(&name))
	require.NotEqual(t, "postgres", name)
	require.NoError(t, pool.Ping(t.Context()))

	var count int
	require.NoError(t, pool.QueryRow(t.Context(), "SELECT count(*) FROM notes").Scan(&count))
	require.Zero(t, count)
}
