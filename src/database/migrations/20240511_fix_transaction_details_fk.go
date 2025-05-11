package migrations

import (
	"app/src/utils"
	"fmt"

	"gorm.io/gorm"
)

// Fix transaction_details foreign key constraint issue
func FixTransactionDetailsForeignKey(db *gorm.DB) error {
	utils.Log.Info("Running migration: Fix transaction_details foreign key constraint")

	// First, identify and delete transaction records with invalid user_subscription_id
	var invalidCount int64
	result := db.Exec(`
		DELETE FROM transaction_details
		WHERE user_subscription_id NOT IN (
			SELECT id FROM user_subscriptions
		)
	`)
	if result.Error != nil {
		return fmt.Errorf("failed to delete invalid transaction records: %w", result.Error)
	}
	invalidCount = result.RowsAffected
	utils.Log.Infof("Deleted %d transaction records with invalid user_subscription_id", invalidCount)

	// Now add the UserSubscription relationship to the model
	err := db.Exec(`
		ALTER TABLE transaction_details
		ADD CONSTRAINT fk_transaction_details_user_subscription 
		FOREIGN KEY (user_subscription_id) 
		REFERENCES user_subscriptions(id)
	`).Error

	if err != nil {
		return fmt.Errorf("failed to add foreign key constraint: %w", err)
	}

	utils.Log.Info("Successfully added foreign key constraint to transaction_details")
	return nil
}
