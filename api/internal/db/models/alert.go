package models

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	ID                   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time `gorm:"index"`
	UserID               uuid.UUID  `gorm:"type:uuid;not null"`
	TokenID              uuid.UUID  `gorm:"type:uuid;not null"`
	Condition            string     `gorm:"not null"`
	NotificationSettings string
	User                 User  `gorm:"foreignKey:UserID"`
	Token                Token `gorm:"foreignKey:TokenID"`
}
