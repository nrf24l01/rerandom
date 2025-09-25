package database

import (
	"fmt"
	"log"

	"github.com/nrf24l01/rerandom/control_tool/core"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RegisterPostgres(cfg *core.Config) (*gorm.DB){
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", 
	cfg.PGHost, cfg.PGUser, cfg.PGPassword, cfg.PGDatabase, cfg.PGPort, cfg.PGSSLMode, cfg.PGTimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to get db instance: %v", err)
	}

	return db
}