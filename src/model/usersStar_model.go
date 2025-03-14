package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersStar struct {
	ID        uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	Stars     int       `gorm:"not null" json:"stars"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
}

func (usersStar *UsersStar) BeforeCreate(_ *gorm.DB) error {
	usersStar.ID = uuid.New()
	return nil
}
