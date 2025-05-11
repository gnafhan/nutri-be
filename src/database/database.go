package database

import (
	"app/src/database/migrations"
	"app/src/database/seeders"
	"app/src/model"
	"log"

	"gorm.io/gorm"
)

func MigrateAndSeed(db *gorm.DB) {
	// Run auto-migrations
	if err := db.AutoMigrate(
		&model.User{},
		&model.Token{},
		&model.Article{},
		&model.ArticleCategory{},
		&model.MealHistory{},
		&model.MealHistoryDetail{},
		&model.ProductToken{},
		&model.Recipe{},
		&model.UsersStar{},
		&model.UsersWeightHeightHistory{},
		&model.UsersWeightHeightTarget{},
		&model.SubscriptionPlan{},
		&model.UserSubscription{},
		&model.TransactionDetail{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	// Run seeders
	seeders.RunSeeder(db)

	// Run custom migrations
	if err := migrations.CreateEnumDay(db); err != nil {
		log.Fatalf("Failed to create enum type: %v", err)
	}
}
