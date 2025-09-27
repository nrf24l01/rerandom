package config

import (
	"os"
	"strconv"
)

type Config struct {
	SpreadsheetID string
	SheetName     string

	REDIS_HOST     string
	REDIS_PASSWORD string
	REDIS_DB       int
	REDIS_KEY      string

	APP_HOST string
}

func BuildConfigFromEnv() *Config {
	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		redisDB = 0 // default to DB 0 if conversion fails
	}
	return &Config{
		SpreadsheetID: os.Getenv("SPREADSHEET_ID"),
		SheetName:     os.Getenv("SHEET_NAME"),
		REDIS_DB:      redisDB,
		REDIS_HOST:    os.Getenv("REDIS_HOST"),
		REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
		REDIS_KEY:     os.Getenv("REDIS_KEY"),
		APP_HOST:      os.Getenv("APP_HOST"),
	}
}