package models

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `gorm:"index"`
	Name             string     `gorm:"not null;unique"`
	CurrentUsdValue  float64    `gorm:"not null;type:numeric"`
	LastUpdated      time.Time  `gorm:"not null;default:now()"`
	Transactions     []Transaction
	HistoricalPrices []HistoricalPrice
	Alerts           []Alert
}
