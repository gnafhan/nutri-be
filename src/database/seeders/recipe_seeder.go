package seeders

import (
	"app/src/model"
	"log"
	"math/rand"

	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedRecipes(db *gorm.DB) {
	var count int64
	db.Model(&model.Recipe{}).Count(&count)

	if count > 0 {
		log.Println("✅ Recipes already seeded, skipping...")
		return
	}

	// Sample markdown content
	instructions := `# Cooking Instructions
1. Heat the pan.
2. Add oil and sauté onions.
3. Add the main ingredients and cook thoroughly.
4. Serve hot and enjoy!`

	ingredients := `# Ingredients
- 1 cup of rice
- 2 eggs
- 1 tbsp soy sauce
- 1 tsp salt`

	labels := []string{"breakfast", "lunch", "dinner"}
	days := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}

	// Sample image placeholders
	dummyImages := []string{
		"https://placehold.co/800x400/png?text=Recipe+Image+1",
		"https://placehold.co/800x400/png?text=Recipe+Image+2",
		"https://placehold.co/800x400/png?text=Recipe+Image+3",
		"https://placehold.co/800x400/png?text=Recipe+Image+4",
		"https://placehold.co/800x400/png?text=Recipe+Image+5",
	}

	var recipes []model.Recipe
	for i := 0; i < 10; i++ {
		imageIndex := rand.Intn(len(dummyImages))
		imageURL := dummyImages[imageIndex]

		recipes = append(recipes, model.Recipe{
			ID:           uuid.New(),
			UserID:       uuid.New(),
			Name:         "Recipe " + fmt.Sprint(i+1),
			Slug:         "recipe-" + fmt.Sprint(i+1),
			Description:  "This is a sample recipe description.",
			Ingredients:  ingredients,
			Instructions: instructions,
			Label:        &labels[rand.Intn(len(labels))],
			Day:          days[rand.Intn(len(days))],
			Image:        &imageURL,
		})
	}

	if err := db.Create(&recipes).Error; err != nil {
		log.Fatalf("Failed to seed recipes: %v", err)
	}

	log.Printf("✅ %d recipes seeded successfully", len(recipes))
}
