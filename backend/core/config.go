package core

import (
	"os"
)

type Config struct {
	APPHost 	  string

	PGHost                  string
	PGPort                  string
	PGUser                  string
	PGPassword              string
	PGDatabase              string
	PGSSLMode               string
	PGTimeZone              string

	AllowOrigins            string

	JWTAccessSecret         string
	JWTRefreshSecret        string

	TestEnv                 bool
	ProductionEnv           bool
}

func BuildConfigFromEnv() (*Config, error) {
	cfg := &Config{
		APPHost:          os.Getenv("APP_HOST"),
		
		AllowOrigins:     os.Getenv("ALLOW_ORIGINS"),

		PGHost:           os.Getenv("POSTGRES_HOST"),
		PGPort:           os.Getenv("POSTGRES_PORT"),
		PGUser:           os.Getenv("POSTGRES_USER"),
		PGPassword:       os.Getenv("POSTGRES_PASSWORD"),
		PGDatabase:       os.Getenv("POSTGRES_DB"),
		PGSSLMode:        os.Getenv("POSTGRES_SSLMODE"),
		PGTimeZone:       os.Getenv("POSTGRES_TIMEZONE"),

		JWTAccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),

		TestEnv:          os.Getenv("TEST_ENV") == "true",
		ProductionEnv:    os.Getenv("PRODUCTION_ENV") == "true",
	}
	return cfg, nil
}
