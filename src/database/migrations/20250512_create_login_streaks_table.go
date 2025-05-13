package migrations

import (
	"app/src/utils"
	"fmt"

	"gorm.io/gorm"
)

// CreateLoginStreaksTable creates the login_streaks table
func CreateLoginStreaksTable(db *gorm.DB) error {
	utils.Log.Info("Running migration: Create login_streaks table")

	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS login_streaks (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			login_date TIMESTAMP NOT NULL,
			day_of_week INTEGER NOT NULL,
			current_streak INTEGER NOT NULL DEFAULT 1,
			longest_streak INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT login_streaks_user_id_idx UNIQUE (user_id, login_date)
		)
	`).Error

	if err != nil {
		return fmt.Errorf("failed to create login_streaks table: %w", err)
	}

	// Create index on user_id for faster lookups
	err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_login_streaks_user_id ON login_streaks(user_id)
	`).Error

	if err != nil {
		return fmt.Errorf("failed to create index on login_streaks: %w", err)
	}

	utils.Log.Info("Successfully created login_streaks table")
	return nil
}
