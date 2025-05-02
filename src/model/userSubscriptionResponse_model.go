package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (userSubscriptionPlanResponse *UserSubscriptionResponse) BeforeCreate(_ *gorm.DB) error {
	userSubscriptionPlanResponse.ID = uuid.New()
	return nil
}
