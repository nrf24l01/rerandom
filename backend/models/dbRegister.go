package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nrf24l01/rerandom/backend/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func RegisterPostgres(cfg *core.Config) (*gorm.DB){
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", 
	cfg.PGHost, cfg.PGUser, cfg.PGPassword, cfg.PGDatabase, cfg.PGPort, cfg.PGSSLMode, cfg.PGTimeZone)
	
	// Настройка логгера GORM в зависимости от окружения
	var gormConfig *gorm.Config
	if os.Getenv("RUNTIME_PRODUCTION") != "true" {
		// В режиме разработки включаем логирование SQL запросов
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Info,   // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      false,         // Include params in SQL log
				Colorful:                  true,          // Colorful SQL logs
			},
		)
		gormConfig = &gorm.Config{
			Logger: newLogger,
		}
	} else {
		// В продакшене используем молчаливый логгер
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("failed to get db instance: %v", err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)

	if err := db.AutoMigrate(&Digit{}, &DigitDrop{}, &User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}