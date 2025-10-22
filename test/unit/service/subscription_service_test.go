package service_test

import (
	"app/src/model"
	"app/src/service"
	"app/test"
	"app/test/helper"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateFreemiumSubscription(t *testing.T) {
	t.Run("should create freemium subscription successfully", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.ClearSubscriptions(test.DB)

		// Create a test user
		userID := uuid.New()
		user := &model.User{
			ID:    userID,
			Name:  "Test User",
			Email: "test@example.com",
		}
		helper.InsertUser(test.DB, user)

		// Create subscription service
		subscriptionService := service.NewSubscriptionService(test.DB, nil)

		// Create freemium subscription
		err := subscriptionService.CreateFreemiumSubscription(nil, userID)
		assert.NoError(t, err)

		// Verify subscription was created
		subscription, err := helper.GetUserSubscription(test.DB, userID)
		assert.NoError(t, err)
		assert.NotNil(t, subscription)
		assert.Equal(t, "freemium_trial", subscription.PaymentMethod)
		assert.Equal(t, "completed", subscription.PaymentStatus)
		assert.True(t, subscription.IsActive)
	})

	t.Run("should not create duplicate freemium subscription", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.ClearSubscriptions(test.DB)

		// Create a test user
		userID := uuid.New()
		user := &model.User{
			ID:    userID,
			Name:  "Test User",
			Email: "test@example.com",
		}
		helper.InsertUser(test.DB, user)

		// Create subscription service
		subscriptionService := service.NewSubscriptionService(test.DB, nil)

		// Create first freemium subscription
		err := subscriptionService.CreateFreemiumSubscription(nil, userID)
		assert.NoError(t, err)

		// Try to create second freemium subscription
		err = subscriptionService.CreateFreemiumSubscription(nil, userID)
		assert.NoError(t, err) // Should not error, just skip

		// Verify only one subscription exists
		subscriptions := []model.UserSubscription{}
		test.DB.Where("user_id = ?", userID).Find(&subscriptions)
		assert.Len(t, subscriptions, 1)
	})
}
