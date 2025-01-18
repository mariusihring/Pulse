package models

import (
	"time"

	"github.com/google/uuid"
)

type Snapshot struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
	SubwalletID  uuid.UUID  `gorm:"type:uuid;not null"`
	SnapshotDate time.Time  `gorm:"not null;default:now()"`
	TotalPnl     float64    `gorm:"not null;type:numeric"`
	Subwallet    Subwallet  `gorm:"foreignKey:SubwalletID;constraint:OnDelete:CASCADE"`
}
