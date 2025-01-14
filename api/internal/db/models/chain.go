package models

import (
	"time"

	"github.com/google/uuid"
)

type Chain struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
	Name       string     `gorm:"not null"`
	Subwallets []Subwallet
}
