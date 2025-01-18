package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenSnapshot struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time     `gorm:"index"`
	SubwalletTokenID uuid.UUID      `gorm:"type:uuid;not null"`
	SnapshotDate     time.Time      `gorm:"not null;default:now()"`
	TotalPnl         float64        `gorm:"not null;type:numeric"`
	SubwalletToken   SubwalletToken `gorm:"foreignKey:SubwalletTokenID;constraint:OnDelete:CASCADE"`
}
