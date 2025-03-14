package validation

import (
	"app/src/model"
	"time"
)

type Register struct {
	Name           string              `json:"name" validate:"required,max=50" example:"fake name"`
	Email          string              `json:"email" validate:"required,email,max=50" example:"fake@example.com"`
	Password       string              `json:"password" validate:"required,min=8,max=20,password" example:"password1"`
	BirthDate      time.Time           `json:"birth_date" validate:"required" example:"2000-01-01T00:00:00Z"`
	Height         float64             `json:"height" validate:"required,gt=0" example:"170.5"`
	Weight         float64             `json:"weight" validate:"required,gt=0" example:"65.5"`
	Gender         model.GenderType    `json:"gender" validate:"required,oneof=Male Female" example:"Male"`
	ActivityLevel  model.ActivityLevel `json:"activity_level" validate:"required,oneof=Light Medium Heavy" example:"Medium"`
	MedicalHistory string              `json:"medical_history" validate:"required,max=1000" example:"No known medical issues"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email,max=50" example:"fake@example.com"`
	Password string `json:"password" validate:"required,min=8,max=20,password" example:"password1"`
}

type GoogleLogin struct {
	Name  string `json:"name" validate:"required,max=50"`
	Email string `json:"email" validate:"required,email,max=50"`
}

type Logout struct {
	RefreshToken string `json:"refresh_token" validate:"required,max=255"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required,max=255"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email,max=50" example:"fake@example.com"`
}

type Token struct {
	Token string `json:"token" validate:"required,max=255"`
}
