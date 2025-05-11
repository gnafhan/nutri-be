package example

// AdminSubscriptionPlanWithUsersResponse is a Swagger-friendly version of response.SubscriptionPlanWithUsers
type AdminSubscriptionPlanWithUsersResponse struct {
	ID             string                     `json:"id"`
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

// AdminSubscriptionPlansResponse is a Swagger-friendly version of response.SuccessWithSubscriptionPlans
type AdminSubscriptionPlansResponse struct {
	Status  string                                   `json:"status"`
	Message string                                   `json:"message"`
	Data    []AdminSubscriptionPlanWithUsersResponse `json:"data"`
}
