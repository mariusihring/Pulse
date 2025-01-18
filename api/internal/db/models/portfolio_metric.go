package models

import (
	"time"

	"github.com/google/uuid"
)

type PortfolioMetric struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null"`
	MetricName   string     `gorm:"not null"`
	MetricValue  float64    `gorm:"not null;type:numeric"`
	CalculatedAt time.Time  `gorm:"not null;default:now()"`
	User         User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
