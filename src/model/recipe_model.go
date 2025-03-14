package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Recipe struct {
	ID           uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID `gorm:"not null" json:"user_id"`
	Name         string    `gorm:"not null" json:"name"`
	Slug         string    `gorm:"unique;not null" json:"slug"`
	Image        *string   `json:"image,omitempty"`
	Description  string    `gorm:"type:text;not null" json:"description"`
	Ingredients  string    `gorm:"type:text;not null" json:"ingredients"`
	Instructions string    `gorm:"type:text;not null" json:"instructions"`
	Label        *string   `json:"label,omitempty"`
	CreatedAt    time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt    time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
}

func (recipe *Recipe) BeforeCreate(_ *gorm.DB) error {
	recipe.ID = uuid.New()
	return nil
}
