package models

import (
	"log"

	"github.com/nrf24l01/go-web-utils/goorm"
	"github.com/nrf24l01/go-web-utils/passhash"
)

type User struct {
	goorm.BaseModel
	Username string `json:"username" gorm:"not null;uniqueIndex"`
	Password string `json:"password" gorm:"not null"`
}

func (u *User) CheckPassword(password string) bool {
	res, err := passhash.CheckPassword(password, u.Password)
	log.Printf("Password check result: %v, error: %v", res, err)
	return res && err == nil
}

func (u *User) SetPassword(password string) error {
	var err error
	u.Password, err = passhash.HashPassword(password, passhash.DefaultParams)
	return err
}
