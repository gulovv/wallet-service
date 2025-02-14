package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	Balance float64   `gorm:"type:decimal(20,2)"`
}