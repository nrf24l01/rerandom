package models

import "github.com/nrf24l01/go-web-utils/goorm"

type User struct {
	goorm.BaseModel
	Username string `json:"username" gorm:"not null;uniqueIndex"`
	Password string `json:"password" gorm:"not null"`
}