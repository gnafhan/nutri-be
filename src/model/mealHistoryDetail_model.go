package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MealHistoryDetail struct {
	ID            uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	MealHistoryID uuid.UUID `gorm:"not null" json:"meal_history_id"`
	APIResult     string    `gorm:"type:text;not null" json:"api_result"`
	CreatedAt     time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt     time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
}

func (mealHistoryDetail *MealHistoryDetail) BeforeCreate(_ *gorm.DB) error {
	mealHistoryDetail.ID = uuid.New()
	return nil
}
