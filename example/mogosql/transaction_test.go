package main_test

import (
	"errors"
	"testing"
	"uuid"

	"github.com/jackc/pgx/v5"
	example "github.com/momoli-dev/mogo/example/mogosql"
	"github.com/momoli-dev/mogo/mogosql"
	"github.com/momoli-dev/mogo/mogosql/mogosqltest"
	"github.com/stretchr/testify/require"
)

func TestWithTx(t *testing.T) {
	t.Parallel()

	const (
		commit   = "commit"
		failure  = "error"
		panicked = "panic"
	)

	for _, outcome := range []string{commit, failure, panicked} {
		t.Run(outcome, func(t *testing.T) {
			t.Parallel()

			pool, _ := mogosqltest.NewPostgres(t, example.MigrationFS)
			ctx := t.Context()
			callbackErr := errors.New("operation failed")
			run := func() error {
				return mogosql.WithTx(ctx, pool, func(tx pgx.Tx) error {
					_, err := tx.Exec(ctx,
						"INSERT INTO notes (note_id, title, body) VALUES ($1, $2, $3)",
						uuid.NewV7().String(), "test", "test",
					)
					if err != nil {
						return err
					}

					switch outcome {
					case failure:
						return callbackErr
					case panicked:
						panic("operation panicked")
					}

					return nil
				})
			}

			switch outcome {
			case commit:
				require.NoError(t, run())
			case failure:
				require.ErrorIs(t, run(), callbackErr)
			case panicked:
				require.PanicsWithValue(t, "operation panicked", func() { _ = run() })
			}

			var count int
			require.NoError(t, pool.QueryRow(ctx, "SELECT count(*) FROM notes").Scan(&count))
			if outcome == commit {
				require.Equal(t, 1, count)
			} else {
				require.Zero(t, count)
			}
		})
	}
}
