package migrations

import (
	"gorm.io/gorm"
)

func CreateEnumDay(db *gorm.DB) error {
	// Create the enum type
	if err := db.Exec(`
		DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'day_enum') THEN
				CREATE TYPE day_enum AS ENUM ('sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday');
			END IF;
		END $$;
	`).Error; err != nil {
		return err
	}

	// Add the column using the enum type
	return db.Exec(`
		ALTER TABLE recipes ADD COLUMN IF NOT EXISTS day day_enum NOT NULL;
	`).Error
}
