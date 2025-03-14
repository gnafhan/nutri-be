package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MealHistory struct {
	ID        uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	Title     string    `gorm:"not null" json:"title"`
	MealTime  time.Time `gorm:"not null" json:"meal_time"`
	Label     *string   `json:"label,omitempty"`
	Calories  float64   `gorm:"type:decimal(6,2);not null" json:"calories"`
	Protein   float64   `gorm:"type:decimal(6,2);not null" json:"protein"`
	Carbs     float64   `gorm:"type:decimal(6,2);not null" json:"carbs"`
	Fat       float64   `gorm:"type:decimal(6,2);not null" json:"fat"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
}

func (mealHistory *MealHistory) BeforeCreate(_ *gorm.DB) error {
	mealHistory.ID = uuid.New()
	return nil
}
