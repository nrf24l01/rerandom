package core

import (
	"os"
	"strconv"
)

type Config struct {
	APPHost 	  string

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
	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		APPHost:          os.Getenv("APP_HOST"),
		
		AllowOrigins:     os.Getenv("ALLOW_ORIGINS"),

		RedisHost: 	os.Getenv("REDIS_HOST"),
		RedisPort: 	os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB: 	redisDB,
		KeysSetName: os.Getenv("KEYS_SET_NAME"),

		TestEnv:          os.Getenv("TEST_ENV") == "true",
		ProductionEnv:    os.Getenv("PRODUCTION_ENV") == "true",
	}
	return cfg, nil
}
