package core

import (
	"os"
)

type Config struct {
	PGHost        string
	PGPort        string
	PGUser        string
	PGPassword    string
	PGDatabase    string
	PGSSLMode     string
	PGTimeZone    string
}

func BuildConfigFromEnv() (*Config) {
	cfg := &Config{
		PGHost:           os.Getenv("PG_HOST"),
		PGPort:           os.Getenv("PG_PORT"),
		PGUser:           os.Getenv("PG_USER"),
		PGPassword:       os.Getenv("PG_PASSWORD"),
		PGDatabase:       os.Getenv("PG_DATABASE"),
		PGSSLMode:        os.Getenv("PG_SSL_MODE"),
		PGTimeZone:       os.Getenv("PG_TIME_ZONE"),
	}
	return cfg
}
