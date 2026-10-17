package dbpool

import (
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func TestConfigRejectsMissingOrUnsafeErrors(t *testing.T) {
	_, err := Config("   ")
	if err == nil || err.Error() != "DATABASE_URL is unset" {
		t.Fatalf("empty url error = %v", err)
	}

	_, err = Config("http://user:super-secret-password@example.com/db")
	if err == nil {
		t.Fatal("expected invalid url")
	}
	if err.Error() != "DATABASE_URL is invalid" {
		t.Fatalf("error = %q", err.Error())
	}
	if strings.Contains(err.Error(), "super-secret-password") {
		t.Fatal("error included the secret")
	}
}

func TestConfigPoolLimits(t *testing.T) {
	cfg, err := Config("postgres://user:super-secret-password@127.0.0.1:5432/db?sslmode=disable")
	if err != nil {
		t.Fatal("expected a parseable url")
	}
	if cfg.MaxConns != 4 || cfg.MinConns != 0 {
		t.Fatalf("conns max=%d min=%d", cfg.MaxConns, cfg.MinConns)
	}
	if cfg.MaxConnIdleTime != 5*time.Minute {
		t.Fatalf("idle = %s", cfg.MaxConnIdleTime)
	}
	if cfg.MaxConnLifetime != 30*time.Minute {
		t.Fatalf("lifetime = %s", cfg.MaxConnLifetime)
	}
	if cfg.HealthCheckPeriod != time.Minute {
		t.Fatalf("health = %s", cfg.HealthCheckPeriod)
	}
	if cfg.PingTimeout != time.Second {
		t.Fatalf("ping timeout = %s", cfg.PingTimeout)
	}
	if cfg.ConnConfig.ConnectTimeout != 5*time.Second {
		t.Fatalf("connect timeout = %s", cfg.ConnConfig.ConnectTimeout)
	}
	if cfg.ConnConfig.DefaultQueryExecMode != pgx.QueryExecModeSimpleProtocol {
		t.Fatalf("exec mode = %v", cfg.ConnConfig.DefaultQueryExecMode)
	}
}
