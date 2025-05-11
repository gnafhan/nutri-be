package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionPlanResponse adalah model untuk response subscription plan
type SubscriptionPlanResponse struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	Price          int             `json:"price"`
	PriceFormatted string          `json:"price_formatted"`
	Description    string          `json:"description"`
	Features       map[string]bool `json:"features"`
	IsRecommended  bool            `json:"is_recommended"`
	ValidityDays   int             `json:"validity_days"`
	AIscanLimit    int             `json:"ai_scan_limit"`
}

// SubscriptionPlanWithUsers adalah model untuk plan dengan users
type SubscriptionPlanWithUsers struct {
	ID             uuid.UUID                  `json:"id"`
	Name           string                     `json:"name"`
	Price          int                        `json:"price"`
	PriceFormatted string                     `json:"price_formatted"`
	Description    string                     `json:"description"`
	AIscanLimit    int                        `json:"ai_scan_limit"`
	ValidityDays   int                        `json:"validity_days"`
	Features       map[string]bool            `json:"features"`
	IsActive       bool                       `json:"is_active"`
	Users          []UserSubscriptionResponse `json:"users,omitempty"`
	UserCount      int                        `json:"user_count"`
}

func (subscriptionPlanResponse *SubscriptionPlanResponse) BeforeCreate(_ *gorm.DB) error {
	subscriptionPlanResponse.ID = uuid.New()
	return nil
}
