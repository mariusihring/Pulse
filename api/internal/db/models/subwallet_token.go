package models

import (
	"time"

	"github.com/google/uuid"
)

type SubwalletToken struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
	SubwalletID uuid.UUID  `gorm:"type:uuid;not null"`
	TokenID     uuid.UUID  `gorm:"type:uuid;not null"`
	Amount      float64    `gorm:"not null;type:numeric"`
	TotalPnl    float64    `gorm:"not null;type:numeric"`
	Subwallet   Subwallet  `gorm:"foreignKey:SubwalletID;constraint:OnDelete:CASCADE"`
	Token       Token      `gorm:"foreignKey:TokenID;constraint:OnDelete:CASCADE"`
	Snapshots   []TokenSnapshot
}
