package core

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	HTTPAddr string
	Postgres PostgresConfig
}

// PostgresConfig describes the shared Postgres server and role. Every module
// connects with the same credentials but to its own database.
type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
	// AdminDB is the maintenance database used to create module databases.
	AdminDB string
}

// LoadConfig reads configuration from environment variables.
// POSTGRES_USER and POSTGRES_PASSWORD are required; everything else has a default.
func LoadConfig() (Config, error) {
	var missing []string
	require := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			missing = append(missing, key)
		}
		return v
	}

	pg := PostgresConfig{
		User:     require("POSTGRES_USER"),
		Password: require("POSTGRES_PASSWORD"),
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		AdminDB:  getEnv("POSTGRES_DB", "postgres"),
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return Config{
		HTTPAddr: ":" + getEnv("HTTP_PORT", "8080"),
		Postgres: pg,
	}, nil
}

// DatabaseName returns the database a module should use: <MODULE>_DB_NAME,
// defaulting to the module name itself (e.g. FINANCE_DB_NAME, default "finance").
func DatabaseName(module string) string {
	return getEnv(strings.ToUpper(module)+"_DB_NAME", module)
}

// URL builds a connection string for the given database on the shared server.
func (p PostgresConfig) URL(database string) string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(p.User, p.Password),
		Host:     net.JoinHostPort(p.Host, p.Port),
		Path:     database,
		RawQuery: url.Values{"sslmode": {p.SSLMode}}.Encode(),
	}
	return u.String()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
