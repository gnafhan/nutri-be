package database

import (
	"app/src/config"
	"app/src/utils"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dbHost, dbName string) *gorm.DB {
	// Validate required database parameters
	if dbHost == "" {
		utils.Log.Errorf("Database host is empty")
		panic("Database host cannot be empty")
	}
	if config.DBUser == "" {
		utils.Log.Errorf("Database user is empty")
		panic("Database user cannot be empty")
	}
	if dbName == "" {
		utils.Log.Errorf("Database name is empty")
		panic("Database name cannot be empty")
	}
	if config.DBPort == 0 {
		utils.Log.Errorf("Database port is 0 or invalid")
		panic("Database port must be a valid port number")
	}

	// hihihi maap
	dsn := fmt.Sprintf(
		"host=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, dbName, config.DBPort,
	)
	
	utils.Log.Infof("Attempting to connect to database with DSN: host=%s user=%s dbname=%s port=%d", 
		dbHost, config.DBUser, dbName, config.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true, // Disable prepared statements to avoid connection issues
		TranslateError:         true,
	})
	if err != nil {
		utils.Log.Errorf("Failed to connect to database: %+v", err)
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
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
