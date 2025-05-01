package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlanResponse struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	Price          int             `json:"price"`
	PriceFormatted string          `json:"price_formatted"`
	Features       map[string]bool `json:"features"`
	IsRecommended  bool            `json:"is_recommended"`
}

func (subscriptionPlanResponse *SubscriptionPlanResponse) BeforeCreate(_ *gorm.DB) error {
	subscriptionPlanResponse.ID = uuid.New()
	return nil
}
