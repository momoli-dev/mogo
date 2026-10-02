// Package mogosql provides utilities for working with SQL databases.
package mogosql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// MigrateParams configures migrations.
type MigrateParams struct {
	Conn         *sql.DB
	MigrationsFS fs.FS
	Directory    string
	Dialect      goose.Dialect
}

// Migrate applies migrations without closing the supplied connection.
func Migrate(ctx context.Context, p MigrateParams) error {
	migrationRoot, err := fs.Sub(p.MigrationsFS, p.Directory)
	if err != nil {
		return fmt.Errorf("open migrations directory: %w", err)
	}

	provider, err := goose.NewProvider(p.Dialect, p.Conn, migrationRoot)
	if err != nil {
		return fmt.Errorf("create migration provider: %w", err)
	}

	res, err := provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	if len(res) == 0 {
		slog.Info("No migrations applied, already up to date")
		return nil
	}

	for _, m := range res {
		slog.Info("Migrations applied", "migration", m.String())
	}

	return nil
}

// PostgresMigrateParams configures PostgreSQL migrations.
type PostgresMigrateParams struct {
	ConnURL      string
	MigrationsFS fs.FS
}

// MigratePostgres applies Goose migrations from migrations/ and closes its connection.
func MigratePostgres(ctx context.Context, p PostgresMigrateParams) (retErr error) {
	db, err := sql.Open("pgx", p.ConnURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		closeErr := db.Close()
		if closeErr != nil {
			retErr = errors.Join(retErr, fmt.Errorf("close database: %w", closeErr))
		}
	}()

	err = Migrate(ctx, MigrateParams{
		Conn:         db,
		MigrationsFS: p.MigrationsFS,
		Dialect:      goose.DialectPostgres,
		Directory:    "migrations",
	})
	if err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	return nil
}
