package database

import (
	"app/src/config"
	"app/src/database/migrations"
	"app/src/database/seeders"
	"app/src/model"
	"app/src/utils"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dbHost, dbName string) *gorm.DB {
	// hihihi maap
	dsn := fmt.Sprintf(
		"host=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, dbName, config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		TranslateError:         true,
	})
	if err != nil {
		utils.Log.Errorf("Failed to connect to database: %+v", err)
	}

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		utils.Log.Errorf("Failed to enable uuid-ossp extension: %+v", err)
	}

	MigrateAndSeed(db)

	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Failed to connect to database: %+v", errDB)
	}

	// Config connection pooling
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	return db
}

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

	// Run product token columns migration (without foreign key constraints)
	if err := db.Exec(`
		ALTER TABLE product_tokens 
		ADD COLUMN IF NOT EXISTS created_by_id UUID DEFAULT NULL,
		ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
	`).Error; err != nil {
		utils.Log.Warnf("Failed to add columns to product_tokens: %v", err)
	}
}
