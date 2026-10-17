// Package config reads process environment for the API.
// The process does not load .env files.
package config

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Desire-in-tech/morrow-and-muse/backend/internal/dbpool"
)

// Config is the process configuration. Do not log it: DatabaseURL is a secret.
type Config struct {
	Port        string
	DatabaseURL string
	CORSOrigins []string
}

// Load reads PORT, DATABASE_URL, and CORS_ALLOWED_ORIGINS.
// PORT defaults to 8080. DATABASE_URL is required. An empty allowlist allows
// no browser origins. Errors never include the database URL or an origin value.
func Load() (Config, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if _, err := dbpool.Config(databaseURL); err != nil {
		return Config{}, err
	}

	port, err := port()
	if err != nil {
		return Config{}, err
	}
	origins, err := origins()
	if err != nil {
		return Config{}, err
	}
	return Config{
		Port:        port,
		DatabaseURL: databaseURL,
		CORSOrigins: origins,
	}, nil
}

func port() (string, error) {
	raw := strings.TrimSpace(os.Getenv("PORT"))
	if raw == "" {
		return "8080", nil
	}
	number, err := strconv.Atoi(raw)
	if err != nil || number < 1 || number > 65535 {
		return "", errors.New("PORT is invalid")
	}
	return strconv.Itoa(number), nil
}

func origins() ([]string, error) {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if raw == "" {
		return nil, nil
	}
	var out []string
	seen := map[string]struct{}{}
	for i, part := range strings.Split(raw, ",") {
		origin := strings.TrimSpace(part)
		if origin == "" {
			return nil, errors.New("CORS_ALLOWED_ORIGINS entry " + strconv.Itoa(i+1) + " is empty")
		}
		if !validOrigin(origin) {
			return nil, errors.New("CORS_ALLOWED_ORIGINS entry " + strconv.Itoa(i+1) + " is invalid")
		}
		if _, ok := seen[origin]; ok {
			continue
		}
		seen[origin] = struct{}{}
		out = append(out, origin)
	}
	return out, nil
}

func validOrigin(origin string) bool {
	if origin == "*" || strings.Contains(origin, "*") {
		return false
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Host == "" || parsed.User != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	if parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.Opaque != "" {
		return false
	}
	return parsed.String() == origin
}
