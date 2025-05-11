package response

import (
	"app/src/model"
)

// SuccessWithPaginateSubscriptions adalah respons untuk daftar subscription dengan pagination
type SuccessWithPaginateSubscriptions struct {
	Status       string                           `json:"status"`
	Message      string                           `json:"message"`
	Results      []model.UserSubscriptionResponse `json:"results"`
	Page         int                              `json:"page"`
	Limit        int                              `json:"limit"`
	TotalPages   int64                            `json:"total_pages"`
	TotalResults int64                            `json:"total_results"`
}

// SuccessWithSubscription adalah respons untuk detail subscription
type SuccessWithSubscription struct {
	Status  string                         `json:"status"`
	Message string                         `json:"message"`
	Data    model.UserSubscriptionResponse `json:"data"`
}

// SubscriptionPlanWithUsers adalah model untuk plan dengan users
type SubscriptionPlanWithUsers struct {
	ID             string                           `json:"id"`
	Name           string                           `json:"name"`
	Price          int                              `json:"price"`
	PriceFormatted string                           `json:"price_formatted"`
	Description    string                           `json:"description"`
	AIscanLimit    int                              `json:"ai_scan_limit"`
	ValidityDays   int                              `json:"validity_days"`
	Features       map[string]bool                  `json:"features"`
	IsActive       bool                             `json:"is_active"`
	Users          []model.UserSubscriptionResponse `json:"users,omitempty"`
	UserCount      int                              `json:"user_count"`
}

// SuccessWithSubscriptionPlans adalah respons untuk daftar subscription plans
type SuccessWithSubscriptionPlans struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []SubscriptionPlanWithUsers `json:"data"`
}

// SuccessWithTransactions is a response for transaction logs
type SuccessWithTransactions struct {
	Status       string                    `json:"status"`
	Message      string                    `json:"message"`
	Data         []model.TransactionDetail `json:"data"`
	Page         int                       `json:"page,omitempty"`
	Limit        int                       `json:"limit,omitempty"`
	TotalPages   int64                     `json:"total_pages,omitempty"`
	TotalResults int64                     `json:"total_results,omitempty"`
}

// SuccessWithTransaction is a response for a single transaction
type SuccessWithTransaction struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Data    model.TransactionDetail `json:"data"`
}
