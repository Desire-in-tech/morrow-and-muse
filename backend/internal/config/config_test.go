package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:super-secret-password@127.0.0.1:5432/db?sslmode=disable")
	t.Setenv("PORT", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", " http://localhost:5173 , http://127.0.0.1:5173 ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Port != "8080" {
		t.Fatalf("port = %q", cfg.Port)
	}
	if len(cfg.CORSOrigins) != 2 || cfg.CORSOrigins[0] != "http://localhost:5173" || cfg.CORSOrigins[1] != "http://127.0.0.1:5173" {
		t.Fatalf("origins = %#v", cfg.CORSOrigins)
	}
	if strings.Contains(cfg.Port, "secret") {
		t.Fatal("port leaked")
	}
}

func TestLoadRejectsMissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "  ")
	t.Setenv("PORT", "8080")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	_, err := Load()
	if err == nil || err.Error() != "DATABASE_URL is unset" {
		t.Fatalf("error = %v", err)
	}
}

func TestLoadHidesInvalidDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "http://user:super-secret-password@example.com/db")
	t.Setenv("PORT", "8080")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	_, err := Load()
	if err == nil || err.Error() != "DATABASE_URL is invalid" {
		t.Fatalf("error = %v", err)
	}
	if strings.Contains(err.Error(), "super-secret-password") || strings.Contains(err.Error(), "example.com") {
		t.Fatal("error leaked the url")
	}
}

func TestLoadRejectsBadPortAndOrigins(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:secret@127.0.0.1:5432/db?sslmode=disable")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")

	t.Setenv("PORT", "0")
	if _, err := Load(); err == nil || err.Error() != "PORT is invalid" {
		t.Fatalf("port error = %v", err)
	}
	t.Setenv("PORT", "abc")
	if _, err := Load(); err == nil || err.Error() != "PORT is invalid" {
		t.Fatalf("port error = %v", err)
	}
	t.Setenv("PORT", "8080")

	for _, origins := range []string{"*", "https://shop.example/*", "http://localhost:5173/shop", "http://localhost:5173,,"} {
		t.Setenv("CORS_ALLOWED_ORIGINS", origins)
		_, err := Load()
		if err == nil || !strings.Contains(err.Error(), "CORS_ALLOWED_ORIGINS") {
			t.Fatalf("origins %q error = %v", origins, err)
		}
		if strings.Contains(err.Error(), "example") || strings.Contains(err.Error(), "localhost") {
			t.Fatalf("error leaked origin %q: %v", origins, err)
		}
	}
}
