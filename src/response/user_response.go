package response

import (
	"time"

	"github.com/google/uuid"
)

type CreateUser struct {
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Role            string     `json:"role"`
	IsEmailVerified bool       `json:"is_email_verified"`
	BirthDate       *time.Time `json:"birth_date,omitempty"`
	Height          *float64   `json:"height,omitempty"`
	Weight          *float64   `json:"weight,omitempty"`
	Gender          *string    `json:"gender,omitempty"`
	ActivityLevel   *string    `json:"activity_level,omitempty"`
	MedicalHistory  *string    `json:"medical_history,omitempty"`
}

type GetUsers struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Role            string     `json:"role"`
	IsEmailVerified bool       `json:"is_email_verified"`
	BirthDate       *time.Time `json:"birth_date,omitempty"`
	Height          *float64   `json:"height,omitempty"`
	Weight          *float64   `json:"weight,omitempty"`
	Gender          *string    `json:"gender,omitempty"`
	ActivityLevel   *string    `json:"activity_level,omitempty"`
	MedicalHistory  *string    `json:"medical_history,omitempty"`
}
