package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LoginStreak represents the user's login streak data
type LoginStreak struct {
	ID            uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	UserID        uuid.UUID `gorm:"not null;index" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"-"`
	LoginDate     time.Time `gorm:"not null" json:"login_date"`
	DayOfWeek     int       `gorm:"not null" json:"day_of_week"` // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	CurrentStreak int       `gorm:"not null;default:1" json:"current_streak"`
	LongestStreak int       `gorm:"not null;default:1" json:"longest_streak"`
	CreatedAt     time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt     time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
}

// LoginStreakResponse represents the response for login streak data
type LoginStreakResponse struct {
	CurrentStreak int                  `json:"current_streak"`
	LongestStreak int                  `json:"longest_streak"`
	WeeklyStreak  []LoginStreakDayInfo `json:"weekly_streak"`
}

// LoginStreakData is an alias for LoginStreakResponse for Swagger documentation
type LoginStreakData = LoginStreakResponse

// LoginStreakDayInfo represents streak info for a specific day of the week
type LoginStreakDayInfo struct {
	DayOfWeek int       `json:"day_of_week"` // 1 = Monday, ..., 7 = Sunday
	Date      time.Time `json:"date"`
	HasLogin  bool      `json:"has_login"`
}

func (streak *LoginStreak) BeforeCreate(_ *gorm.DB) error {
	streak.ID = uuid.New() // Generate UUID before create
	return nil
}
