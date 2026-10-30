package mogosql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithTx commits on success and rolls back on error or panic.
func WithTx(ctx context.Context, pool *pgxpool.Pool, fn func(pgx.Tx) error) (retErr error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to pool.Begin: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			retErr = errors.Join(retErr, fmt.Errorf("failed to tx.Rollback: %w", err))
		}
	}()

	if err := fn(tx); err != nil {
		return fmt.Errorf("failed to fn: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to tx.Commit: %w", err)
	}

	return nil
}
