package validation

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionQuery adalah struktur untuk query parameter subscription
type SubscriptionQuery struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Status string `query:"status"`
}

// UpdateSubscription adalah struktur untuk update subscription
type UpdateSubscription struct {
	PlanID        *uuid.UUID `json:"plan_id" validate:"omitempty,uuid"`
	IsActive      *bool      `json:"is_active" validate:"omitempty"`
	AIscansUsed   *int       `json:"ai_scans_used" validate:"omitempty,min=0"`
	StartDate     *time.Time `json:"start_date" validate:"omitempty"`
	EndDate       *time.Time `json:"end_date" validate:"omitempty"`
	PaymentMethod *string    `json:"payment_method" validate:"omitempty"`
}

// UpdatePaymentStatus adalah struktur untuk update payment status
type UpdatePaymentStatus struct {
	Status string `json:"status" validate:"required,oneof=pending success failed"`
}

// UpdateSubscriptionPlan adalah struktur untuk update subscription plan
type UpdateSubscriptionPlan struct {
	Name         *string          `json:"name" validate:"omitempty,min=2,max=50"`
	Price        *int             `json:"price" validate:"omitempty,min=1"`
	Description  *string          `json:"description" validate:"omitempty"`
	AIscanLimit  *int             `json:"ai_scan_limit" validate:"omitempty,min=1"`
	ValidityDays *int             `json:"validity_days" validate:"omitempty,min=1"`
	Features     *map[string]bool `json:"features" validate:"omitempty"`
	IsActive     *bool            `json:"is_active" validate:"omitempty"`
}
