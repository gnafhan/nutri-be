package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlan struct {
	ID           uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()"`
	Name         string    `gorm:"not null"`
	Price        int       `gorm:"not null"` // in Rupiah
	Description  string
	AIscanLimit  int       `gorm:"not null"` // -1 for unlimited
	ValidityDays int       `gorm:"not null"` // in days
	Features     string    `gorm:"type:jsonb"`
	IsActive     bool      `gorm:"default:true"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (subscriptionPlan *SubscriptionPlan) BeforeCreate(_ *gorm.DB) error {
	subscriptionPlan.ID = uuid.New()
	return nil
}
