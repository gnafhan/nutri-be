package service

import (
	"app/src/model"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginStreakService interface {
	RecordLogin(userID uuid.UUID) error
	GetLoginStreak(userID uuid.UUID) (*model.LoginStreakResponse, error)
}

type loginStreakServiceImpl struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewLoginStreakService(db *gorm.DB, validate *validator.Validate) LoginStreakService {
	return &loginStreakServiceImpl{
		DB:       db,
		Validate: validate,
	}
}

func (service *loginStreakServiceImpl) RecordLogin(userID uuid.UUID) error {
	// Check if user exists
	var user model.User
	if err := service.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	// Check if user already logged in today
	var existingStreak model.LoginStreak
	err := service.DB.Where("user_id = ? AND login_date >= ?", userID, todayStart).First(&existingStreak).Error
	if err == nil {
		// User already logged in today, no need to update
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other error occurred
		return err
	}

	// Get the most recent login streak record
	var latestStreak model.LoginStreak
	err = service.DB.Where("user_id = ?", userID).Order("login_date DESC").First(&latestStreak).Error

	// Calculate current streak
	currentStreak := 1
	longestStreak := 1

	if err == nil {
		// Check if the last login was yesterday
		yesterday := todayStart.AddDate(0, 0, -1)
		if latestStreak.LoginDate.Year() == yesterday.Year() &&
			latestStreak.LoginDate.Month() == yesterday.Month() &&
			latestStreak.LoginDate.Day() == yesterday.Day() {
			// Consecutive login, increment streak
			currentStreak = latestStreak.CurrentStreak + 1
		}

		// Update longest streak if needed
		longestStreak = latestStreak.LongestStreak
		if currentStreak > longestStreak {
			longestStreak = currentStreak
		}
	}

	// Create new streak record
	newStreak := model.LoginStreak{
		UserID:        userID,
		LoginDate:     todayStart,
		DayOfWeek:     int(today.Weekday()),
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
	}

	return service.DB.Create(&newStreak).Error
}

func (service *loginStreakServiceImpl) GetLoginStreak(userID uuid.UUID) (*model.LoginStreakResponse, error) {
	// Check if user exists
	var user model.User
	if err := service.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Get the most recent login streak record to get current and longest streak
	var latestStreak model.LoginStreak
	err := service.DB.Where("user_id = ?", userID).Order("login_date DESC").First(&latestStreak).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User has no login streaks yet
			return &model.LoginStreakResponse{
				CurrentStreak: 0,
				LongestStreak: 0,
				WeeklyStreak:  getEmptyWeeklyStreak(),
			}, nil
		}
		return nil, err
	}

	// Get weekly streak data
	today := time.Now()
	weekStart := getStartOfWeek(today)
	weekEnd := weekStart.AddDate(0, 0, 6)

	var weeklyStreaks []model.LoginStreak
	err = service.DB.Where("user_id = ? AND login_date BETWEEN ? AND ?",
		userID, weekStart, weekEnd).
		Order("login_date ASC").
		Find(&weeklyStreaks).Error
	if err != nil {
		return nil, err
	}

	// Format weekly streak data
	weeklyStreakInfo := formatWeeklyStreakData(weekStart, weeklyStreaks)

	return &model.LoginStreakResponse{
		CurrentStreak: latestStreak.CurrentStreak,
		LongestStreak: latestStreak.LongestStreak,
		WeeklyStreak:  weeklyStreakInfo,
	}, nil
}

// Helper function to get the start of the week (Monday)
func getStartOfWeek(date time.Time) time.Time {
	offset := int(date.Weekday())
	if offset == 0 {
		offset = 7 // Sunday is 0, but we want it to be 7 for calculation
	}
	return time.Date(date.Year(), date.Month(), date.Day()-offset+1, 0, 0, 0, 0, date.Location())
}

// Helper function to format weekly streak data
func formatWeeklyStreakData(weekStart time.Time, streaks []model.LoginStreak) []model.LoginStreakDayInfo {
	weeklyData := make([]model.LoginStreakDayInfo, 7)

	// Initialize with all days of the week
	for i := 0; i < 7; i++ {
		date := weekStart.AddDate(0, 0, i)
		weeklyData[i] = model.LoginStreakDayInfo{
			DayOfWeek: i + 1, // 1 = Monday, ..., 7 = Sunday
			Date:      date,
			HasLogin:  false,
		}
	}

	// Mark days with login
	for _, streak := range streaks {
		// Convert from Go's weekday (0=Sunday) to our format (1=Monday, 7=Sunday)
		dayIndex := streak.DayOfWeek
		if dayIndex == 0 {
			dayIndex = 6 // Sunday becomes index 6 (7th day)
		} else {
			dayIndex -= 1 // Other days shift by 1
		}

		if dayIndex >= 0 && dayIndex < 7 {
			weeklyData[dayIndex].HasLogin = true
		}
	}

	return weeklyData
}

// Helper function to get empty weekly streak data
func getEmptyWeeklyStreak() []model.LoginStreakDayInfo {
	weeklyData := make([]model.LoginStreakDayInfo, 7)
	weekStart := getStartOfWeek(time.Now())

	for i := 0; i < 7; i++ {
		date := weekStart.AddDate(0, 0, i)
		weeklyData[i] = model.LoginStreakDayInfo{
			DayOfWeek: i + 1, // 1 = Monday, ..., 7 = Sunday
			Date:      date,
			HasLogin:  false,
		}
	}

	return weeklyData
}
