package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
	Name      string     `gorm:"not null"`
	Email     string     `gorm:"unique;not null"`
	Password  string     `gorm:"not null"`
	Wallets   []Wallet
	Alerts    []Alert
	Metrics   []PortfolioMetric
	Roles     []Role `gorm:"many2many:user_roles"`
}
