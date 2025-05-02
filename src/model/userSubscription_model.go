package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSubscription struct {
	ID            uuid.UUID        `gorm:"primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID        `gorm:"not null"`
	PlanID        uuid.UUID        `gorm:"not null"`
	Plan          SubscriptionPlan `gorm:"foreignKey:PlanID"`
	AIscansUsed   int              `gorm:"default:0"`
	StartDate     time.Time        `gorm:"not null"`
	EndDate       time.Time        `gorm:"not null"`
	IsActive      bool             `gorm:"default:true"`
	PaymentMethod string           `gorm:"size:50"`
	TransactionID string           `gorm:"size:100"`
	CreatedAt     time.Time        `gorm:"autoCreateTime"`
}

type PurchaseSubscriptionRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=gopay shopeepay bank_transfer"`
}

func (userSubscription *UserSubscription) BeforeCreate(_ *gorm.DB) error {
	userSubscription.ID = uuid.New()
	return nil
}
