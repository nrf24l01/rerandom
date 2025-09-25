package models

import (
	"github.com/google/uuid"
	"github.com/nrf24l01/go-web-utils/goorm"
)

type Digit struct {
	goorm.BaseModel
	Value     int    `json:"answ" gorm:"not null"` // Переименовал для соответствия API (answ)
	Max       *int   `json:"max,omitempty" gorm:"default:null"`
	Min       *int   `json:"min,omitempty" gorm:"default:null"`
	MaxDrops  int    `json:"max_drops" gorm:"not null;default:1"`
}

type DigitDrop struct {
	goorm.BaseModel
	DigitID   uuid.UUID `json:"digit_id" gorm:"type:uuid;not null;index"`
	Digit     Digit     `json:"digit" gorm:"foreignKey:DigitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}