package example

import (
	"time"

	"github.com/google/uuid"
)

// SimpleUser adalah contoh user sederhana untuk Swagger
type SimpleUser struct {
	ID            uuid.UUID `json:"id" example:"e088d183-9eea-4a11-8d5d-74d7ec91bdf5"`
	Name          string    `json:"name" example:"John Doe"`
	Email         string    `json:"email" example:"john.doe@example.com"`
	Role          string    `json:"role" example:"user"`
	VerifiedEmail bool      `json:"verified_email" example:"true"`
}

// ProductTokenResponse adalah contoh untuk product token
type ProductTokenResponse struct {
	ID          uuid.UUID   `json:"id" example:"e088d183-9eea-4a11-8d5d-74d7ec91bdf5"`
	UserID      uuid.UUID   `json:"user_id" example:"a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6"`
	User        *SimpleUser `json:"user,omitempty"`
	Token       string      `json:"token" example:"abc123def456ghi789"`
	CreatedByID uuid.UUID   `json:"created_by_id" example:"b1c2d3e4-f5g6-h7i8-j9k0-l1m2n3o4p5q6"`
	CreatedBy   *SimpleUser `json:"created_by,omitempty"`
	ActivatedAt *time.Time  `json:"activated_at,omitempty" example:"2025-04-20T14:30:00Z"`
	IsActive    bool        `json:"is_active" example:"true"`
	CreatedAt   time.Time   `json:"created_at" example:"2025-04-01T10:00:00Z"`
	UpdatedAt   time.Time   `json:"updated_at" example:"2025-04-01T10:00:00Z"`
}

// GetAllProductTokensResponse adalah contoh untuk respons get all product tokens
type GetAllProductTokensResponse struct {
	Status  string                 `json:"status" example:"success"`
	Message string                 `json:"message" example:"Product tokens retrieved successfully"`
	Data    []ProductTokenResponse `json:"data"`
}

// CreateProductTokenRequest adalah contoh untuk request pembuatan custom token
type CreateProductTokenRequest struct {
	Token    string `json:"token" example:"custom123token"`
	IsActive bool   `json:"is_active" example:"true"`
}

// CreateProductTokenResponse adalah contoh untuk create product token
type CreateProductTokenResponse struct {
	Status  string               `json:"status" example:"success"`
	Message string               `json:"message" example:"Product token created successfully"`
	Data    ProductTokenResponse `json:"data"`
}

// DeleteProductTokenResponse adalah contoh untuk delete product token
type DeleteProductTokenResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Product token deleted successfully"`
}

// SuccessWithPaginateUsers adalah contoh untuk respons pagination users
type SuccessWithPaginateUsers struct {
	Status       string       `json:"status" example:"success"`
	Message      string       `json:"message" example:"Get all users successfully"`
	Results      []SimpleUser `json:"results"`
	Page         int          `json:"page" example:"1"`
	Limit        int          `json:"limit" example:"10"`
	TotalPages   int64        `json:"total_pages" example:"5"`
	TotalResults int64        `json:"total_results" example:"42"`
}
