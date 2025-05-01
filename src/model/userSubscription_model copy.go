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
type UserSubscriptionResponse struct {
	ID            uuid.UUID                `json:"id"`
	UserID        uuid.UUID                `json:"user_id"`
	Plan          SubscriptionPlanResponse `json:"plan"`
	AIscansUsed   int                      `json:"ai_scans_used"`
	StartDate     time.Time                `json:"start_date"`
	EndDate       time.Time                `json:"end_date"`
	IsActive      bool                     `json:"is_active"`
	PaymentMethod string                   `json:"payment_method"`
	CreatedAt     time.Time                `json:"created_at"`
}

func (userSubscription *UserSubscription) BeforeCreate(_ *gorm.DB) error {
	userSubscription.ID = uuid.New()
	return nil
}
