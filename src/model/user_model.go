package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenderType string
type ActivityLevel string

const (
	Male   GenderType = "Male"
	Female GenderType = "Female"

	Light  ActivityLevel = "Light"
	Medium ActivityLevel = "Medium"
	Heavy  ActivityLevel = "Heavy"
)

type User struct {
	ID             uuid.UUID      `gorm:"primaryKey;not null" json:"id"`
	Name           string         `gorm:"not null" json:"name"`
	Email          string         `gorm:"uniqueIndex;not null" json:"email"`
	Password       string         `gorm:"not null" json:"-"`
	Role           string         `gorm:"default:user;not null" json:"role"`
	VerifiedEmail  bool           `gorm:"default:false;not null" json:"verified_email"`
	ProfilePicture string         `gorm:"default:null" json:"profile_picture"`
	GoogleIDToken  string         `gorm:"default:null" json:"google_id_token"`
	Phone          string         `gorm:"size:20;default:null" json:"phone"`
	BirthDate      *time.Time     `gorm:"default:null" json:"birth_date"`
	Height         *float64       `gorm:"type:decimal(5,2);default:null" json:"height"`
	Weight         *float64       `gorm:"type:decimal(5,2);default:null" json:"weight"`
	Gender         *GenderType    `gorm:"type:varchar(10);default:null" json:"gender"`
	ActivityLevel  *ActivityLevel `gorm:"type:varchar(10);default:null" json:"activity_level"`
	MedicalHistory *string        `gorm:"type:text;default:null" json:"medical_history"`
	CreatedAt      time.Time      `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt      time.Time      `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
	Token          []Token        `gorm:"foreignKey:user_id;references:id" json:"-"`
}

func (user *User) BeforeCreate(_ *gorm.DB) error {
	user.ID = uuid.New() // Generate UUID before create
	return nil
}
