package validation

// ProductToken adalah struktur untuk validasi product token
type ProductToken struct {
	Token string `json:"token" validate:"required"`
}

// CreateCustomToken adalah struktur untuk validasi custom token
type CreateCustomToken struct {
	Token    string `json:"token" validate:"required,min=8,max=32"`
	IsActive bool   `json:"is_active" validate:"boolean"`
}

// ProductTokenQuery adalah struktur untuk query parameter product token
type ProductTokenQuery struct {
	WithUser bool `query:"with_user"`
}
