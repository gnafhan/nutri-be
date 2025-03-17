package validation

import (
	"app/src/model"
	"time"
)

type CreateUser struct {
	Name     string `json:"name" validate:"required,max=50" example:"fake name"`
	Email    string `json:"email" validate:"required,email,max=50" example:"fake@example.com"`
	Password string `json:"password" validate:"required,min=8,max=20,password" example:"password1"`
	Role     string `json:"role" validate:"required,oneof=user admin,max=50" example:"user"`
}

type UpdateUser struct {
	Name           string               `json:"name,omitempty" validate:"omitempty,max=50" example:"fake name"`
	Email          string               `json:"email,omitempty" validate:"omitempty,email,max=50" example:"fake@example.com"`
	Password       string               `json:"password,omitempty" validate:"omitempty,min=8,max=20,password" example:"password1"`
	ProfilePicture string               `json:"profile_picture,omitempty" validate:"omitempty,url,max=255" example:"https://example.com/profile.jpg"`
	BirthDate      *time.Time           `json:"birth_date,omitempty" validate:"omitempty"`
	Height         *float64             `json:"height,omitempty" validate:"omitempty,gte=0,lte=300" example:"175.5"`
	Weight         *float64             `json:"weight,omitempty" validate:"omitempty,gte=0,lte=500" example:"70.3"`
	Gender         *model.GenderType    `json:"gender,omitempty" validate:"omitempty,oneof=Male Female" example:"Male"`
	ActivityLevel  *model.ActivityLevel `json:"activity_level,omitempty" validate:"omitempty,oneof=Light Medium Heavy" example:"Medium"`
	MedicalHistory *string              `json:"medical_history,omitempty" validate:"omitempty,max=1000" example:"No known allergies"`
}

type UpdatePassOrVerify struct {
	Password      string `json:"password,omitempty" validate:"omitempty,min=8,max=20,password" example:"password1"`
	VerifiedEmail bool   `json:"verified_email" swaggerignore:"true" validate:"omitempty,boolean"`
}

type QueryUser struct {
	Page   int    `validate:"omitempty,number,max=50"`
	Limit  int    `validate:"omitempty,number,max=50"`
	Search string `validate:"omitempty,max=50"`
}
