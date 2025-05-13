package response

import (
	"app/src/model"
	"time"
)

// SuccessWithLoginStreak represents the response for login streak data
type SuccessWithLoginStreak struct {
	Status  string                `json:"status" example:"success"`
	Message string                `json:"message" example:"Login streak retrieved successfully"`
	Data    model.LoginStreakData `json:"data"`
}

// LoginStreakData represents the login streak data structure
type LoginStreakData struct {
	CurrentStreak int                  `json:"current_streak" example:"5"`
	LongestStreak int                  `json:"longest_streak" example:"12"`
	WeeklyStreak  []LoginStreakDayInfo `json:"weekly_streak"`
}

// LoginStreakDayInfo represents streak info for a specific day of the week
type LoginStreakDayInfo struct {
	DayOfWeek int       `json:"day_of_week" example:"1"` // 1 = Monday, ..., 7 = Sunday
	Date      time.Time `json:"date" example:"2023-06-05T00:00:00Z"`
	HasLogin  bool      `json:"has_login" example:"true"`
}

// Example responses for Swagger documentation
type LoginStreakExample struct {
	// Success response example
	SuccessResponse SuccessWithLoginStreak
}

// GetSuccessExample returns an example of a successful login streak response
func (e *LoginStreakExample) GetSuccessExample() SuccessWithLoginStreak {
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1) // Get Monday of current week

	weeklyStreak := make([]model.LoginStreakDayInfo, 7)
	for i := 0; i < 7; i++ {
		date := weekStart.AddDate(0, 0, i)
		hasLogin := i < 5 // Monday to Friday has login
		weeklyStreak[i] = model.LoginStreakDayInfo{
			DayOfWeek: i + 1, // 1 = Monday, ..., 7 = Sunday
			Date:      date,
			HasLogin:  hasLogin,
		}
	}

	return SuccessWithLoginStreak{
		Status:  "success",
		Message: "Login streak retrieved successfully",
		Data: model.LoginStreakData{
			CurrentStreak: 5,
			LongestStreak: 12,
			WeeklyStreak:  weeklyStreak,
		},
	}
}
