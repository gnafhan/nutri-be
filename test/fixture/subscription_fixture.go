package fixture

import (
	"app/src/model"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// FreemiumPlan returns a Freemium Trial subscription plan fixture
func FreemiumPlan() *model.SubscriptionPlan {
	features := map[string]bool{
		"scan_ai":         true,
		"scan_calorie":    true,
		"chatbot":         true,
		"bmi_check":       true,
		"weight_tracking": true,
		"health_info":     true,
	}
	featuresJSON, _ := json.Marshal(features)

	return &model.SubscriptionPlan{
		ID:           uuid.New(),
		Name:         "Freemium Trial",
		Price:        0,
		Description:  "14-day free trial with full features",
		AIscanLimit:  10,
		ValidityDays: 14,
		Features:     string(featuresJSON),
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
}

// ActiveFreemiumSubscription returns an active freemium subscription fixture
func ActiveFreemiumSubscription(userID uuid.UUID, planID uuid.UUID) *model.UserSubscription {
	now := time.Now()
	endDate := now.AddDate(0, 0, 14) // 14 days from now

	return &model.UserSubscription{
		ID:            uuid.New(),
		UserID:        userID,
		PlanID:        planID,
		StartDate:     now,
		EndDate:       endDate,
		IsActive:      true,
		PaymentMethod: "freemium_trial",
		PaymentStatus: "completed",
		AIscansUsed:   0,
		CreatedAt:     now,
	}
}

// ExpiredFreemiumSubscription returns an expired freemium subscription fixture
func ExpiredFreemiumSubscription(userID uuid.UUID, planID uuid.UUID) *model.UserSubscription {
	now := time.Now()
	startDate := now.AddDate(0, 0, -15) // Started 15 days ago
	endDate := now.AddDate(0, 0, -1)    // Expired 1 day ago

	return &model.UserSubscription{
		ID:            uuid.New(),
		UserID:        userID,
		PlanID:        planID,
		StartDate:     startDate,
		EndDate:       endDate,
		IsActive:      false,
		PaymentMethod: "freemium_trial",
		PaymentStatus: "completed",
		AIscansUsed:   10,
		CreatedAt:     startDate,
	}
}

// PaidSubscription returns a paid subscription fixture
func PaidSubscription(userID uuid.UUID, planID uuid.UUID) *model.UserSubscription {
	now := time.Now()
	endDate := now.AddDate(0, 0, 30) // 30 days from now

	return &model.UserSubscription{
		ID:            uuid.New(),
		UserID:        userID,
		PlanID:        planID,
		StartDate:     now,
		EndDate:       endDate,
		IsActive:      true,
		PaymentMethod: "credit_card",
		PaymentStatus: "completed",
		AIscansUsed:   5,
		CreatedAt:     now,
	}
}

// UserWithFreemium returns a user with active freemium subscription fixture
func UserWithFreemium() *model.User {
	return &model.User{
		ID:            uuid.New(),
		Name:          "Freemium User",
		Email:         "freemium@example.com",
		Password:      "hashedpassword123",
		Role:          "user",
		VerifiedEmail: true,
		CreatedAt:     time.Now(),
	}
}

// UserWithExpiredFreemium returns a user with expired freemium subscription fixture
func UserWithExpiredFreemium() *model.User {
	return &model.User{
		ID:            uuid.New(),
		Name:          "Expired Freemium User",
		Email:         "expired@example.com",
		Password:      "hashedpassword123",
		Role:          "user",
		VerifiedEmail: true,
		CreatedAt:     time.Now(),
	}
}

// UserWithPaidSubscription returns a user with paid subscription fixture
func UserWithPaidSubscription() *model.User {
	return &model.User{
		ID:            uuid.New(),
		Name:          "Paid User",
		Email:         "paid@example.com",
		Password:      "hashedpassword123",
		Role:          "user",
		VerifiedEmail: true,
		CreatedAt:     time.Now(),
	}
}
