// Package db applies the embedded goose migrations.
//
// goose is pinned in go.mod (github.com/pressly/goose/v3). Run it through
// this package so tests and the command use the same SQL:
//
//	DATABASE_URL=postgres://... go run ./cmd/migrate up
//	DATABASE_URL=postgres://... go run ./cmd/migrate down-to 0
package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Up applies every migration, including the generated seed.
func Up(ctx context.Context, databaseURL string) error {
	provider, err := newProvider(databaseURL)
	if err != nil {
		return err
	}
	defer provider.Close()
	_, err = provider.Up(ctx)
	return err
}

// DownTo rolls back through the given goose version. Version 0 removes all
// migrations in this module.
func DownTo(ctx context.Context, databaseURL string, version int64) error {
	provider, err := newProvider(databaseURL)
	if err != nil {
		return err
	}
	defer provider.Close()
	_, err = provider.DownTo(ctx, version)
	return err
}

func newProvider(databaseURL string) (*goose.Provider, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("database URL is empty")
	}
	migrations, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return nil, fmt.Errorf("migration files: %w", err)
	}
	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, migrations)
	if err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("goose provider: %w", err)
	}
	return provider, nil
}
