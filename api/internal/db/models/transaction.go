package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time          `gorm:"index"`
	TokenID         uuid.UUID           `gorm:"type:uuid;not null"`
	TransactionType string              `gorm:"not null"`
	Amount          float64             `gorm:"not null;type:numeric"`
	TransactionDate time.Time           `gorm:"not null;default:now()"`
	CategoryID      uuid.UUID           `gorm:"type:uuid;not null"`
	Token           Token               `gorm:"foreignKey:TokenID"`
	Category        TransactionCategory `gorm:"foreignKey:CategoryID"`
}
