package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductToken struct {
	ID          uuid.UUID  `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID  `gorm:"default:null" json:"user_id"`
	Token       string     `gorm:"unique;not null" json:"token"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt   time.Time  `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
}

func (productToken *ProductToken) BeforeCreate(_ *gorm.DB) error {
	productToken.ID = uuid.New()
	return nil
}
