package seeders

import (
	"app/src/model"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedSubscriptionPlans(db *gorm.DB) {
	var count int64
	db.Model(&model.SubscriptionPlan{}).Count(&count)

	if count > 0 {
		log.Println("✅ Subscription plans already seeded, skipping...")
		return
	}

	plans := []model.SubscriptionPlan{
		createPlan(
			"Hemat",
			15000,
			10,
			30,
			map[string]bool{
				"scan_ai":         true,
				"scan_calorie":    true,
				"chatbot":         true,
				"bmi_check":       false,
				"weight_tracking": false,
				"health_info":     false,
			},
			"Paket dasar untuk pemula",
		),
		createPlan(
			"Early Bird",
			99000,
			60,
			90,
			map[string]bool{
				"scan_ai":         true,
				"scan_calorie":    true,
				"chatbot":         true,
				"bmi_check":       true,
				"weight_tracking": true,
				"health_info":     true,
			},
			"Paket premium dengan semua fitur",
			true, // Mark as best seller
		),
		createPlan(
			"Sehat",
			30000,
			10,
			30,
			map[string]bool{
				"scan_ai":         true,
				"scan_calorie":    true,
				"chatbot":         true,
				"bmi_check":       true,
				"weight_tracking": false,
				"health_info":     false,
			},
			"Paket best seller dengan fitur lengkap",
		),
		createPlan(
			"Sultan",
			120000,
			30,
			90,
			map[string]bool{
				"scan_ai":         true,
				"scan_calorie":    true,
				"chatbot":         true,
				"bmi_check":       true,
				"weight_tracking": true,
				"health_info":     true,
			},
			"Paket premium dengan semua fitur",
		),
	}

	if err := db.Create(&plans).Error; err != nil {
		log.Fatalf("Failed to seed subscription plans: %v", err)
	}

	log.Printf("✅ %d subscription plans seeded successfully", len(plans))
}

// Helper function to create plan with JSON features
func createPlan(
	name string,
	price int,
	scanLimit int,
	validityDays int,
	features map[string]bool,
	description string,
	isBestSeller ...bool,
) model.SubscriptionPlan {
	featuresJSON, _ := json.Marshal(features)

	plan := model.SubscriptionPlan{
		ID:           uuid.New(),
		Name:         name,
		Price:        price,
		Description:  description,
		AIscanLimit:  scanLimit,
		ValidityDays: validityDays,
		Features:     string(featuresJSON),
		IsActive:     true,
	}

	if len(isBestSeller) > 0 && isBestSeller[0] {
		plan.Description += " (Best Seller)"
	}

	return plan
}
