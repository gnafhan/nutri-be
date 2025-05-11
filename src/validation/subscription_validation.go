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
