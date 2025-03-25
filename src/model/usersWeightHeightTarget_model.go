package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersWeightHeightTarget struct {
	ID         uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID     uuid.UUID `gorm:"not null" json:"user_id"`
	Weight     float64   `gorm:"type:decimal(5,2);not null" json:"weight"`
	Height     float64   `gorm:"type:decimal(5,2);not null" json:"height"`
	TargetDate time.Time `gorm:"autoCreateTime" json:"target_date"`
	CreatedAt  time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt  time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
}

func (usersWeightHeightTarget *UsersWeightHeightTarget) BeforeCreate(_ *gorm.DB) error {
	usersWeightHeightTarget.ID = uuid.New()
	return nil
}
