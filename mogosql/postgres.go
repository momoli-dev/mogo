package mogosql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	pgxgeom "github.com/twpayne/pgx-geom"
)

// PostgresPoolParams configures a PostgreSQL connection pool.
type PostgresPoolParams struct {
	ConnStr    string
	HasPostgis bool
}

// NewPostgresPool opens and pings a pool. The caller owns its cleanup.
func NewPostgresPool(ctx context.Context, p PostgresPoolParams) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(p.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to pgxpool.ParseConfig: %w", err)
	}

	if p.HasPostgis {
		config.AfterConnect = pgxgeom.Register
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to pgxpool.NewWithConfig: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("failed to pool.Ping: %w", err)
	}

	return pool, nil
}
