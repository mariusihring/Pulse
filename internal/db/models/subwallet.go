package models

import (
	"time"

	"github.com/google/uuid"
)

type Subwallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
	Name      string     `gorm:"not null"`
	WalletID  uuid.UUID  `gorm:"type:uuid;not null"`
	ChainID   uuid.UUID  `gorm:"type:uuid;not null"`
	Wallet    Wallet     `gorm:"foreignKey:WalletID"`
	Chain     Chain      `gorm:"foreignKey:ChainID"`
	Tokens    []SubwalletToken
	Snapshots []Snapshot
	Address   string
}
