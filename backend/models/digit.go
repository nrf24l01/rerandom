package models

import "github.com/nrf24l01/go-web-utils/goorm"

type Digit struct {
	goorm.BaseModel
	Value int `json:"value" gorm:"not null"`
	IfMax int `json:"if_max" gorm:"not null;default:0"`
	IfMin int `json:"if_min" gorm:"not null;default:0"`
	DropCount int `json:"drop_count" gorm:"not null;default:0"`
}

type DigitDrop struct {
	goorm.BaseModel
	DigitID string `json:"digit_id" gorm:"not null;index"`
	Digit   Digit  `json:"digit" gorm:"foreignKey:DigitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}