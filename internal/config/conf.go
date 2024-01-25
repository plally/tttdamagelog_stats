package config

import (
	"os"
)

type Config struct {
	BasicAuth   map[string]string
	PostgresURL string
}

func Get() *Config {
	authMap := make(map[string]string)
	if os.Getenv("BASIC_AUTH_USERNAME") != "" && os.Getenv("BASIC_AUTH_PASSWORD") != "" {
		authMap[os.Getenv("BASIC_AUTH_USERNAME")] = os.Getenv("BASIC_AUTH_PASSWORD")
	}

	return &Config{
		BasicAuth:   authMap,
		PostgresURL: os.Getenv("POSTGRES_URL"),
	}
}
