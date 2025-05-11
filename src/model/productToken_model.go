package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductToken struct {
	ID          uuid.UUID  `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID  `gorm:"default:null" json:"user_id"`
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Token       string     `gorm:"unique;not null" json:"token"`
	CreatedByID uuid.UUID  `gorm:"default:null" json:"created_by_id"`
	CreatedBy   *User      `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"updated_at"`
}

func (productToken *ProductToken) BeforeCreate(_ *gorm.DB) error {
	productToken.ID = uuid.New()
	return nil
}
