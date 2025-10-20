package integration

import (
	"app/src/model"
	"app/src/response"
	"app/src/validation"
	"app/test"
	"app/test/fixture"
	"app/test/helper"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFreemiumFlow(t *testing.T) {
	t.Run("Complete freemium registration and verification flow", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.ClearSubscriptions(test.DB)

		// Step 1: Register new user
		var requestBody = validation.Register{
			Name:           "Freemium User",
			Email:          "freemium@example.com",
			Password:       "password1",
			BirthDate:      time.Now().AddDate(-25, 0, 0),
			Height:         175.0,
			Weight:         70.0,
			Gender:         "Male",
			ActivityLevel:  "Medium",
			MedicalHistory: "No known allergies",
		}

		bodyJSON, err := json.Marshal(requestBody)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(string(bodyJSON)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")

		apiResponse, err := test.App.Test(request)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, apiResponse.StatusCode)

		// Parse response to get user ID
		bytes, err := io.ReadAll(apiResponse.Body)
		assert.Nil(t, err)

		responseBody := new(response.SuccessWithTokens)
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		userID := responseBody.User.ID
		assert.Equal(t, false, responseBody.User.VerifiedEmail) // Should be unverified initially

		// Step 2: Get verification token from database and verify email (this should create freemium subscription)
		// The verification token is stored in the database, not returned in the registration response
		var verificationToken model.Token
		err = test.DB.Where("user_id = ? AND type = ?", userID, "verify_email").First(&verificationToken).Error
		assert.Nil(t, err)
		assert.NotEmpty(t, verificationToken.Token)

		verifyRequest := httptest.NewRequest(http.MethodGet, "/v1/auth/verify-email?token="+verificationToken.Token, nil)

		verifyResponse, err := test.App.Test(verifyRequest)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, verifyResponse.StatusCode)

		// Step 3: Check that freemium subscription was created
		subscription, err := helper.GetUserSubscription(test.DB, userID)
		assert.Nil(t, err)
		assert.NotNil(t, subscription)
		assert.Equal(t, "freemium_trial", subscription.PaymentMethod)
		assert.Equal(t, "completed", subscription.PaymentStatus)
		assert.True(t, subscription.IsActive)
		assert.True(t, time.Now().Before(subscription.EndDate))

		// Step 4: Access protected endpoint with freemium subscription
		accessToken := responseBody.Tokens.Access.Token
		mealRequest := httptest.NewRequest(http.MethodGet, "/v1/meals", nil)
		mealRequest.Header.Set("Authorization", "Bearer "+accessToken)

		mealResponse, err := test.App.Test(mealRequest)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, mealResponse.StatusCode)
	})

	t.Run("Freemium expiry scenario", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.ClearSubscriptions(test.DB)

		// Create user with expired freemium subscription
		user := fixture.UserWithExpiredFreemium()
		helper.InsertUser(test.DB, user)

		// Create expired freemium subscription
		err := helper.CreateExpiredFreemiumSubscription(test.DB, user.ID)
		assert.Nil(t, err)

		// Generate access token for user
		accessToken, err := fixture.AccessToken(user)
		assert.Nil(t, err)

		// Try to access protected endpoint with expired freemium
		mealRequest := httptest.NewRequest(http.MethodGet, "/v1/meals", nil)
		mealRequest.Header.Set("Authorization", "Bearer "+accessToken)

		mealResponse, err := test.App.Test(mealRequest)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusForbidden, mealResponse.StatusCode)

		// Verify response contains upgrade message
		bytes, err := io.ReadAll(mealResponse.Body)
		assert.Nil(t, err)

		var errorResponse response.ErrorResponse
		err = json.Unmarshal(bytes, &errorResponse)
		assert.Nil(t, err)
		assert.Equal(t, "freemium_expired", errorResponse.Message)
	})

	t.Run("User with existing subscription should not get freemium", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.ClearSubscriptions(test.DB)

		// Create user with existing paid subscription
		user := fixture.UserWithPaidSubscription()
		helper.InsertUser(test.DB, user)

		// Create paid subscription
		var freemiumPlan model.SubscriptionPlan
		test.DB.First(&freemiumPlan, "name = ?", "Freemium Trial")

		paidPlan := &model.SubscriptionPlan{
			ID:           fixture.FreemiumPlan().ID,
			Name:         "Paid Plan",
			Price:        50000,
			Description:  "Paid subscription",
			AIscanLimit:  100,
			ValidityDays: 30,
			Features:     fixture.FreemiumPlan().Features,
			IsActive:     true,
		}
		test.DB.Create(paidPlan)

		paidSubscription := &model.UserSubscription{
			UserID:        user.ID,
			PlanID:        paidPlan.ID,
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 30),
			IsActive:      true,
			PaymentMethod: "credit_card",
			PaymentStatus: "completed",
			AIscansUsed:   0,
		}
		test.DB.Create(paidSubscription)

		// Generate access token for user
		accessToken, err := fixture.AccessToken(user)
		assert.Nil(t, err)

		// Try to access protected endpoint - should work with paid subscription
		mealRequest := httptest.NewRequest(http.MethodGet, "/v1/meals", nil)
		mealRequest.Header.Set("Authorization", "Bearer "+accessToken)

		mealResponse, err := test.App.Test(mealRequest)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, mealResponse.StatusCode)

		// Verify user still has only one subscription (the paid one)
		var subscriptionCount int64
		test.DB.Model(&model.UserSubscription{}).Where("user_id = ?", user.ID).Count(&subscriptionCount)
		assert.Equal(t, int64(1), subscriptionCount)
	})
}
