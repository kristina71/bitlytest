package config

import "os"

type Cfg struct {
	DbDsn     string
	DbDialect string
}

func New() Cfg {
	return Cfg{
		DbDsn:     readFromEnv("DB_DSN", "postgresql://postgres:postgres@localhost:5432/bitlytest?sslmode=disable"),
		DbDialect: readFromEnv("DB_DIALECT", "postgres"),
	}
}

func readFromEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return value
}
