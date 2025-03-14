package seeders

import (
	"app/src/model"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 6)
	for i := range token {
		token[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(token)
}

func generateUniqueToken(db *gorm.DB, existingTokens map[string]bool) string {
	var token string
	for {
		token = generateRandomToken()
		if !existingTokens[token] {
			var count int64
			db.Model(&model.ProductToken{}).Where("token = ?", token).Count(&count)
			if count == 0 {
				existingTokens[token] = true
				break
			}
		}
	}
	return token
}

func SeedProductTokens(db *gorm.DB) {
	var count int64
	db.Model(&model.ProductToken{}).Count(&count)

	if count > 0 {
		log.Println("✅ Product tokens already seeded, skipping...")
		return
	}

	productTokens := []model.ProductToken{}
	existingTokens := make(map[string]bool)

	for i := 0; i < 10; i++ {
		productTokens = append(productTokens, model.ProductToken{
			ID:          uuid.New(),
			UserID:      uuid.Nil, // NULL
			Token:       generateUniqueToken(db, existingTokens),
			ActivatedAt: nil, // NULL
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	if err := db.Create(&productTokens).Error; err != nil {
		log.Fatalf("Failed to seed product tokens: %v", err)
	}

	log.Println("✅ Product tokens seeding completed!")
}

func RunSeeder(db *gorm.DB) {
	log.Println("🚀 Running database seeder...")
	SeedProductTokens(db)
	log.Println("🎉 Database seeding completed!")
}
