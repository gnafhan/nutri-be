package database

import (
	"app/src/database/migrations"
	"app/src/database/seeders"
	"app/src/model"
	"app/src/utils"
	"log"

	"gorm.io/gorm"
)

func MigrateAndSeed(db *gorm.DB) {
	// Run our custom migrations first to fix any data issues
	runCustomMigrations(db)

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
		&model.LoginStreak{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	// Run seeders
	seeders.RunSeeder(db)
}

func runCustomMigrations(db *gorm.DB) {
	// Fix product token records with invalid user_id
	if err := db.Exec(`
		UPDATE product_tokens 
		SET user_id = NULL 
		WHERE user_id IS NOT NULL AND NOT EXISTS (
			SELECT 1 FROM users WHERE users.id = product_tokens.user_id
		);
	`).Error; err != nil {
		utils.Log.Warnf("Failed to clean invalid user_id in product_tokens: %v", err)
	}

	// Run custom enum day migration
	if err := migrations.CreateEnumDay(db); err != nil {
		log.Fatalf("Failed to create enum type: %v", err)
	}

	// Fix transaction details foreign key constraint
	if err := migrations.FixTransactionDetailsForeignKey(db); err != nil {
		utils.Log.Warnf("Failed to fix transaction details foreign key: %v", err)
	}

	// Create login streaks table
	if err := migrations.CreateLoginStreaksTable(db); err != nil {
		utils.Log.Warnf("Failed to create login streaks table: %v", err)
	}

	// Run product token columns migration (without foreign key constraints)
	if err := db.Exec(`
		ALTER TABLE product_tokens 
		ADD COLUMN IF NOT EXISTS created_by_id UUID DEFAULT NULL,
		ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
	`).Error; err != nil {
		utils.Log.Warnf("Failed to add columns to product_tokens: %v", err)
	}
}
