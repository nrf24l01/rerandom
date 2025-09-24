package core

import (
	"os"
)

type Config struct {
	APPHost 	  string

	PGHost        string
	PGPort        string
	PGUser        string
	PGPassword    string
	PGDatabase    string
	PGSSLMode     string
	PGTimeZone    string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	KeysSetName   string

	AllowOrigins            string

	TestEnv                 bool
	ProductionEnv           bool
}

func BuildConfigFromEnv() (*Config, error) {
	cfg := &Config{
		APPHost:          os.Getenv("APP_HOST"),
		
		AllowOrigins:     os.Getenv("ALLOW_ORIGINS"),
		PGHost:           os.Getenv("PG_HOST"),
		PGPort:           os.Getenv("PG_PORT"),
		PGUser:           os.Getenv("PG_USER"),
		PGPassword:       os.Getenv("PG_PASSWORD"),
		PGDatabase:       os.Getenv("PG_DATABASE"),
		PGSSLMode:        os.Getenv("PG_SSL_MODE"),
		PGTimeZone:       os.Getenv("PG_TIME_ZONE"),

		TestEnv:          os.Getenv("TEST_ENV") == "true",
		ProductionEnv:    os.Getenv("PRODUCTION_ENV") == "true",
	}
	return cfg, nil
}
