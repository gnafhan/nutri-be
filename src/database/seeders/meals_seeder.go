package seeders

import (
	"app/src/model"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedMeals(db *gorm.DB) {
	var count int64
	db.Model(&model.MealHistory{}).Count(&count)

	if count > 0 {
		log.Println("✅ Meals already seeded, skipping...")
		return
	}

	// Fixed user ID
	userID, _ := uuid.Parse("7ea4f83c-64c4-4079-ad44-2c3d1e06be2d")

	// Meal types and their typical macro distributions
	mealTypes := []struct {
		title    string
		calories [2]int // min, max range
		protein  [2]int
		carbs    [2]int
		fat      [2]int
		images   []string
	}{
		{
			title:    "Breakfast",
			calories: [2]int{300, 600},
			protein:  [2]int{15, 30},
			carbs:    [2]int{30, 60},
			fat:      [2]int{10, 25},
			images:   []string{"https://example.com/breakfast1.jpg", "https://example.com/breakfast2.jpg", "https://example.com/breakfast3.jpg"},
		},
		{
			title:    "Lunch",
			calories: [2]int{500, 800},
			protein:  [2]int{25, 45},
			carbs:    [2]int{45, 80},
			fat:      [2]int{20, 35},
			images:   []string{"https://example.com/lunch1.jpg", "https://example.com/lunch2.jpg", "https://example.com/lunch3.jpg"},
		},
		{
			title:    "Dinner",
			calories: [2]int{600, 900},
			protein:  [2]int{30, 50},
			carbs:    [2]int{40, 70},
			fat:      [2]int{25, 40},
			images:   []string{"https://example.com/dinner1.jpg", "https://example.com/dinner2.jpg", "https://example.com/dinner3.jpg"},
		},
		{
			title:    "Snack",
			calories: [2]int{100, 300},
			protein:  [2]int{5, 15},
			carbs:    [2]int{15, 30},
			fat:      [2]int{5, 15},
			images:   []string{"https://example.com/snack1.jpg", "https://example.com/snack2.jpg", "https://example.com/snack3.jpg"},
		},
	}

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	var meals []model.MealHistory
	now := time.Now()

	// Generate at least one meal per day for the past month
	for i := 0; i < 30; i++ {
		// Calculate the date (i days ago)
		mealDate := now.Add(-time.Hour * 24 * time.Duration(i))

		// Random number of meals per day (1-3)
		numMeals := rand.Intn(3) + 1

		// Create meals for this day
		for j := 0; j < numMeals; j++ {
			// Select random meal type
			mealType := mealTypes[rand.Intn(len(mealTypes))]

			// Calculate random meal time for this day
			hours := 8 + rand.Intn(13) // Between 8 AM and 8 PM
			mealTime := time.Date(
				mealDate.Year(), mealDate.Month(), mealDate.Day(),
				hours, rand.Intn(60), 0, 0, mealDate.Location(),
			)

			// Generate random values within the specified ranges
			calories := rand.Intn(mealType.calories[1]-mealType.calories[0]+1) + mealType.calories[0]
			protein := rand.Intn(mealType.protein[1]-mealType.protein[0]+1) + mealType.protein[0]
			carbs := rand.Intn(mealType.carbs[1]-mealType.carbs[0]+1) + mealType.carbs[0]
			fat := rand.Intn(mealType.fat[1]-mealType.fat[0]+1) + mealType.fat[0]

			// Select a random image for this meal type
			mealImage := mealType.images[rand.Intn(len(mealType.images))]

			meal := model.MealHistory{
				ID:             uuid.New(),
				UserID:         userID,
				Title:          mealType.title,
				MealTime:       mealTime,
				Calories:       float64(calories),
				Protein:        float64(protein),
				Carbs:          float64(carbs),
				Fat:            float64(fat),
				MealImage:      mealImage,
				Recommendation: nil,
				CreatedAt:      mealTime,
				UpdatedAt:      mealTime,
			}

			meals = append(meals, meal)
		}
	}

	if err := db.Create(&meals).Error; err != nil {
		log.Fatalf("Failed to seed meals: %v", err)
	}

	log.Printf("✅ %d meals seeded successfully for the past month", len(meals))
}
