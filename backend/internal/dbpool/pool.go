// Package dbpool opens the pgx pool used by the API.
//
// The API does not migrate. Apply schema changes with cmd/migrate.
package dbpool

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// MaxConns stays small so one free-tier service does not exhaust
	// Supabase's connection limit.
	maxConns          = 4
	minConns          = 0
	maxConnIdleTime   = 5 * time.Minute
	maxConnLifetime   = 30 * time.Minute
	healthCheckPeriod = time.Minute
	connectTimeout    = 5 * time.Second
	pingTimeout       = time.Second
)

// Config parses databaseURL and applies the pool limits. It does not connect.
// Errors never include the URL.
func Config(databaseURL string) (*pgxpool.Config, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, errors.New("DATABASE_URL is unset")
	}
	cfg, err := pgxpool.ParseConfig(strings.TrimSpace(databaseURL))
	if err != nil {
		return nil, errors.New("DATABASE_URL is invalid")
	}
	cfg.MaxConns = maxConns
	cfg.MinConns = minConns
	cfg.MaxConnIdleTime = maxConnIdleTime
	cfg.MaxConnLifetime = maxConnLifetime
	cfg.HealthCheckPeriod = healthCheckPeriod
	cfg.PingTimeout = pingTimeout
	cfg.ConnConfig.ConnectTimeout = connectTimeout
	// Simple protocol avoids prepared statements, so the same binary works
	// with Supabase's transaction-mode pooler and with a direct connection.
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	return cfg, nil
}

// New opens a pool. With MinConns at 0 it does not dial until the first query.
func New(databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := Config(databaseURL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, errors.New("database pool could not be created")
	}
	return pool, nil
}
