package mogosql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	pgxgeom "github.com/twpayne/pgx-geom"
)

// PostgresDB provides utilities for working with PostgreSQL databases using `pgx` driver.
type PostgresDB struct {
	pool  *pgxpool.Pool
	sqlDB *sql.DB
}

// PostgresDBParams to configure a connection to DB.
type PostgresDBParams struct {
	ConnStr    string
	HasPostgis bool
}

// NewPostgresDB connects to PostgreSQL database using connStr. It opens and manages a connection pool using `pgx`.
func NewPostgresDB(ctx context.Context, p PostgresDBParams) (*PostgresDB, error) {
	config, err := pgxpool.ParseConfig(p.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("postgres parse ConnStr: %w", err)
	}

	if p.HasPostgis {
		config.AfterConnect = pgxgeom.Register
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("postgres new with config: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}

	return &PostgresDB{pool: pool, sqlDB: stdlib.OpenDBFromPool(pool)}, nil
}

// Pool returns the underlying pool implementation from `pgx`.
func (db *PostgresDB) Pool() *pgxpool.Pool {
	return db.pool
}

// Conn returns a standard [sql.Conn] backed by underlying pool. Callers need to close it.
func (db *PostgresDB) Conn(ctx context.Context) (*sql.Conn, error) {
	conn, err := db.sqlDB.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire postgres connection: %w", err)
	}
	return conn, nil
}

// Close cleans up resources managed by [PostgresDB].
func (db *PostgresDB) Close() error {
	err := db.sqlDB.Close()
	db.pool.Close()
	if err != nil {
		return fmt.Errorf("close postgres SQL handle: %w", err)
	}

	return nil
}

// Ping checks the database connection.
func (db *PostgresDB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// WithTx executes begins a new transaction and executed fn within it. Roll back and commit is handled automatically.
func (db *PostgresDB) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	return runTx(ctx, tx, fn)
}

func runTx(ctx context.Context, tx pgx.Tx, fn func(pgx.Tx) error) (retErr error) {
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			retErr = errors.Join(retErr, fmt.Errorf("rollback transaction: %w", err))
		}
	}()

	if err := fn(tx); err != nil {
		return fmt.Errorf("query transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
