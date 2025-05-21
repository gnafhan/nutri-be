package validation

// ProductToken adalah struktur untuk validasi product token
type ProductToken struct {
	Token string `json:"token" validate:"required"`
}

// CreateCustomToken adalah struktur untuk validasi custom token
type CreateCustomToken struct {
	Token              string  `json:"token" validate:"required,min=8,max=32"`
	IsActive           bool    `json:"is_active" validate:"boolean"`
	SubscriptionPlanID *string `json:"subscription_plan_id,omitempty" validate:"omitempty,uuid4"`
}

// ProductTokenQuery adalah struktur untuk query parameter product token
type ProductTokenQuery struct {
	WithUser bool `query:"with_user"`
}

// UpdateProductToken adalah struktur untuk validasi pembaruan product token
type UpdateProductToken struct {
	Token              *string `json:"token,omitempty" validate:"omitempty,min=8,max=32"`
	IsActive           *bool   `json:"is_active,omitempty" validate:"omitempty,boolean"`
	SubscriptionPlanID *string `json:"subscription_plan_id,omitempty" validate:"omitempty,uuid4"` // Allow empty string to clear the plan
}
