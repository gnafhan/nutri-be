package example

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID  `json:"id" example:"e088d183-9eea-4a11-8d5d-74d7ec91bdf5"`
	Name           string     `json:"name" example:"fake name"`
	Email          string     `json:"email" example:"fake@example.com"`
	Role           string     `json:"role" example:"user"`
	VerifiedEmail  bool       `json:"verified_email" example:"false"`
	BirthDate      *time.Time `json:"birth_date,omitempty" example:"2000-01-01T00:00:00Z"`
	Height         *float64   `json:"height,omitempty" example:"175.5"`
	Weight         *float64   `json:"weight,omitempty" example:"65.2"`
	Gender         *string    `json:"gender,omitempty" example:"Male"`
	ActivityLevel  *string    `json:"activity_level,omitempty" example:"Medium"`
	MedicalHistory *string    `json:"medical_history,omitempty" example:"No known allergies"`
}

type GoogleUser struct {
	ID             uuid.UUID  `json:"id" example:"e088d183-9eea-4a11-8d5d-74d7ec91bdf5"`
	Name           string     `json:"name" example:"fake name"`
	Email          string     `json:"email" example:"fake@example.com"`
	Role           string     `json:"role" example:"user"`
	VerifiedEmail  bool       `json:"verified_email" example:"true"`
	BirthDate      *time.Time `json:"birth_date,omitempty" example:"2000-01-01T00:00:00Z"`
	Height         *float64   `json:"height,omitempty" example:"175.5"`
	Weight         *float64   `json:"weight,omitempty" example:"65.2"`
	Gender         *string    `json:"gender,omitempty" example:"Male"`
	ActivityLevel  *string    `json:"activity_level,omitempty" example:"Medium"`
	MedicalHistory *string    `json:"medical_history,omitempty" example:"No known allergies"`
}
