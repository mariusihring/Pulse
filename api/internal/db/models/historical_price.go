package models

import (
	"time"

	"github.com/google/uuid"
)

type HistoricalPrice struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
	TokenID   uuid.UUID  `gorm:"type:uuid;not null"`
	Date      time.Time  `gorm:"not null"`
	Price     float64    `gorm:"not null;type:numeric"`
	Token     Token      `gorm:"foreignKey:TokenID;constraint:OnDelete:CASCADE"`
}
