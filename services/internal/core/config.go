package core

import (
	"fmt"
	"net"
	"net/url"
	"os"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
}

// LoadConfig reads configuration from environment variables.
// POSTGRES_USER, POSTGRES_PASSWORD and POSTGRES_DB are required;
// POSTGRES_HOST, POSTGRES_PORT, POSTGRES_SSLMODE and HTTP_PORT have defaults.
func LoadConfig() (Config, error) {
	var missing []string
	require := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			missing = append(missing, key)
		}
		return v
	}

	user := require("POSTGRES_USER")
	password := require("POSTGRES_PASSWORD")
	name := require("POSTGRES_DB")
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %v", missing)
	}

	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, password),
		Host:     net.JoinHostPort(getEnv("POSTGRES_HOST", "localhost"), getEnv("POSTGRES_PORT", "5432")),
		Path:     name,
		RawQuery: url.Values{"sslmode": {getEnv("POSTGRES_SSLMODE", "disable")}}.Encode(),
	}

	return Config{
		DatabaseURL: dsn.String(),
		HTTPAddr:    ":" + getEnv("HTTP_PORT", "8080"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
