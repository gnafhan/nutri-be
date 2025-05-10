package service

import (
	"app/src/model"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentGateway interface {
	Charge(amount int, method string) (*PaymentResponse, error)
	Refund(transactionID string) error
	CreateTransaction(orderID string, amount int, userDetails map[string]interface{}, paymentMethod string) (*PaymentToken, error)
	CheckTransactionStatus(transactionID string) (interface{}, error)
	HandleNotification(notificationJSON []byte) (interface{}, error)
}

type PaymentResponse struct {
	TransactionID string
	Status        string
}

type SubscriptionService interface {
	GetAllPlans(ctx *fiber.Ctx) ([]model.SubscriptionPlanResponse, error)
	PurchasePlan(ctx *fiber.Ctx, userID uuid.UUID, planID uuid.UUID, paymentMethod string) (*model.PaymentResponse, error)
	GetUserActiveSubscription(ctx *fiber.Ctx, userID uuid.UUID) (*model.UserSubscriptionResponse, error)
	CheckFeatureAccess(ctx *fiber.Ctx, userID uuid.UUID, feature string) (bool, error)
	IncrementScanUsage(ctx *fiber.Ctx, userID uuid.UUID) error
	GetRemainingScans(ctx *fiber.Ctx, userID uuid.UUID) (int, error)
	HandlePaymentNotification(ctx *fiber.Ctx, notificationData []byte) error
}

type subscriptionService struct {
	DB      *gorm.DB
	Log     *logrus.Logger
	Payment PaymentGateway
}

func formatCurrency(amount int) string {
	return fmt.Sprintf("Rp %d", amount)
}

func NewSubscriptionService(db *gorm.DB, payment PaymentGateway) SubscriptionService {
	return &subscriptionService{
		DB:      db,
		Log:     logrus.New(),
		Payment: payment,
	}
}

func (s *subscriptionService) GetAllPlans(ctx *fiber.Ctx) ([]model.SubscriptionPlanResponse, error) {
	var plans []model.SubscriptionPlan
	if err := s.DB.WithContext(ctx.Context()).
		Where("is_active = ?", true).
		Find(&plans).Error; err != nil {
		return nil, err
	}

	var responses []model.SubscriptionPlanResponse
	for _, plan := range plans {
		var features map[string]bool
		if err := json.Unmarshal([]byte(plan.Features), &features); err != nil {
			return nil, err
		}

		responses = append(responses, model.SubscriptionPlanResponse{
			ID:             plan.ID,
			Name:           plan.Name,
			Price:          plan.Price,
			PriceFormatted: formatCurrency(plan.Price),
			Features:       features,
			IsRecommended:  plan.Name == "Early Bird",
			Description:    plan.Description,
			ValidityDays:   plan.ValidityDays,
			AIscanLimit:    plan.AIscanLimit,
		})
	}

	return responses, nil
}

func (s *subscriptionService) PurchasePlan(ctx *fiber.Ctx, userID uuid.UUID, planID uuid.UUID, paymentMethod string) (*model.PaymentResponse, error) {
	var plan model.SubscriptionPlan
	if err := s.DB.WithContext(ctx.Context()).First(&plan, "id = ?", planID).Error; err != nil {
		return nil, errors.New("subscription plan not found")
	}

	fmt.Println("paymentMethod", paymentMethod)
	// Get user details
	var user model.User
	if err := s.DB.WithContext(ctx.Context()).First(&user, "id = ?", userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Generate a unique order ID
	orderID := fmt.Sprintf("SUB-%s-%d", userID.String()[:8], time.Now().Unix())

	// Create a new subscription with pending status
	subscription := model.UserSubscription{
		UserID:        userID,
		PlanID:        planID,
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 0, plan.ValidityDays),
		PaymentMethod: paymentMethod,
		TransactionID: orderID,
		PaymentStatus: "pending",
		IsActive:      false, // Will be activated after payment is completed
	}

	// Save subscription to database
	if err := s.DB.WithContext(ctx.Context()).Create(&subscription).Error; err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Prepare user details for Midtrans
	userDetails := map[string]interface{}{
		"first_name": user.Name,
		"last_name":  "",
		"email":      user.Email,
		"phone":      user.Phone,
	}

	// Create transaction in Midtrans
	paymentToken, err := s.Payment.CreateTransaction(orderID, plan.Price, userDetails, paymentMethod)
	if err != nil {
		// Rollback subscription creation if payment fails
		s.DB.WithContext(ctx.Context()).Delete(&subscription)
		return nil, fmt.Errorf("payment creation failed: %w", err)
	}

	// Return payment details
	return &model.PaymentResponse{
		TransactionToken: paymentToken.Token,
		RedirectURL:      paymentToken.RedirectURL,
		OrderID:          orderID,
	}, nil
}

func (s *subscriptionService) HandlePaymentNotification(ctx *fiber.Ctx, notificationData []byte) error {
	// Log raw notification data
	s.Log.Infof("Processing raw notification data: %s", string(notificationData))

	// Convert to a map for easier access and to ensure we have the data even if Midtrans API call fails
	var notification map[string]interface{}
	if err := json.Unmarshal(notificationData, &notification); err != nil {
		s.Log.Errorf("Failed to parse notification JSON: %v", err)
		return fmt.Errorf("failed to parse notification: %w", err)
	}

	// Extract necessary fields
	orderID, ok := notification["order_id"].(string)
	if !ok {
		s.Log.Error("Order ID missing from notification")
		return errors.New("order_id missing from notification")
	}

	s.Log.Infof("Processing notification for order ID: %s", orderID)

	// Extract transaction status (with fallback)
	transactionStatusStr := "unknown"
	if status, ok := notification["transaction_status"].(string); ok {
		transactionStatusStr = status
	}
	s.Log.Infof("Transaction status: %s", transactionStatusStr)

	// Try to get transaction status from Midtrans but don't fail if it doesn't work
	var transactionStatus interface{}
	transactionStatus, err := s.Payment.HandleNotification(notificationData)
	if err != nil {
		// Check if the error is related to signature verification
		if err.Error() == "invalid signature key" || err.Error() == "error verifying signature" {
			s.Log.Error("Signature verification failed - Potential security threat. Rejecting notification.")
			return fmt.Errorf("security verification failed: %w", err)
		}

		s.Log.Warnf("Could not verify with Midtrans API, proceeding with notification data: %v", err)
		// We'll continue with just the notification data for other types of errors
	} else {
		// Log transaction status response for debugging
		s.Log.Infof("Transaction status from Midtrans: %v", transactionStatus)
	}

	// Find subscription in database
	var subscription model.UserSubscription
	if err := s.DB.WithContext(ctx.Context()).
		Where("transaction_id = ?", orderID).
		First(&subscription).Error; err != nil {
		s.Log.Errorf("Subscription not found for order ID %s: %v", orderID, err)
		return fmt.Errorf("subscription not found with order ID %s: %w", orderID, err)
	}

	s.Log.Infof("Found subscription: ID=%s, UserID=%s, Status=%s",
		subscription.ID, subscription.UserID, subscription.PaymentStatus)

	// Process based on transaction status
	switch transactionStatusStr {
	case "capture", "settlement":
		// Payment success
		s.Log.Infof("Updating subscription %s to success status", subscription.ID)
		subscription.PaymentStatus = "success"
		subscription.IsActive = true
	case "deny", "cancel", "expire":
		// Payment failed
		s.Log.Infof("Updating subscription %s to failed status", subscription.ID)
		subscription.PaymentStatus = "failed"
		subscription.IsActive = false
	case "pending":
		// Payment pending, no changes needed
		s.Log.Infof("Subscription %s remains in pending status", subscription.ID)
	default:
		s.Log.Warnf("Unhandled transaction status for subscription %s: %s",
			subscription.ID, transactionStatusStr)
	}

	// Update the subscription
	if err := s.DB.WithContext(ctx.Context()).Save(&subscription).Error; err != nil {
		s.Log.Errorf("Failed to update subscription %s: %v", subscription.ID, err)
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	// Save detailed transaction information
	transactionDetail := s.createTransactionDetailFromNotification(subscription.ID, notification, notificationData)
	if err := s.DB.WithContext(ctx.Context()).Create(&transactionDetail).Error; err != nil {
		s.Log.Errorf("Failed to save transaction details: %v", err)
		// Continue even if saving details fails
	} else {
		s.Log.Infof("Saved transaction details with ID: %s", transactionDetail.ID)
	}

	s.Log.Infof("Successfully updated subscription %s to status: %s",
		subscription.ID, subscription.PaymentStatus)
	return nil
}

// createTransactionDetailFromNotification creates a TransactionDetail object from the notification data
func (s *subscriptionService) createTransactionDetailFromNotification(
	subscriptionID uuid.UUID,
	notification map[string]interface{},
	rawData []byte,
) *model.TransactionDetail {
	// Create new transaction detail
	detail := &model.TransactionDetail{
		UserSubscriptionID: subscriptionID,
		RawResponse:        model.JSON(rawData),
	}

	// Function to safely get string value from notification
	getStringValue := func(key string) *string {
		if value, ok := notification[key].(string); ok {
			return &value
		}
		return nil
	}

	// Fill common fields
	detail.OrderID = getString(notification, "order_id", "")
	detail.TransactionID = getString(notification, "transaction_id", "")
	detail.TransactionStatus = getString(notification, "transaction_status", "")
	detail.StatusCode = getString(notification, "status_code", "")
	detail.StatusMessage = getString(notification, "status_message", "")
	detail.PaymentType = getString(notification, "payment_type", "")
	detail.GrossAmount = getString(notification, "gross_amount", "")
	detail.Currency = getString(notification, "currency", "")
	detail.FraudStatus = getString(notification, "fraud_status", "")

	// Parse transaction time
	if txTime, ok := notification["transaction_time"].(string); ok {
		if parsedTime, err := time.Parse("2006-01-02 15:04:05", txTime); err == nil {
			detail.TransactionTime = parsedTime
		} else {
			s.Log.Warnf("Failed to parse transaction time: %v", err)
			detail.TransactionTime = time.Now()
		}
	} else {
		detail.TransactionTime = time.Now()
	}

	// Parse settlement time if present
	if settlementTime, ok := notification["settlement_time"].(string); ok {
		if parsedTime, err := time.Parse("2006-01-02 15:04:05", settlementTime); err == nil {
			detail.SettlementTime = &parsedTime
		}
	}

	// Credit card specific fields
	detail.MaskedCard = getStringValue("masked_card")
	detail.CardType = getStringValue("card_type")
	detail.Bank = getStringValue("bank")
	detail.ApprovalCode = getStringValue("approval_code")
	detail.ECI = getStringValue("eci")
	detail.ChannelResponseCode = getStringValue("channel_response_code")
	detail.ChannelResponseMessage = getStringValue("channel_response_message")

	// Bank transfer specific fields
	detail.PermataVANumber = getStringValue("permata_va_number")
	detail.BillerCode = getStringValue("biller_code")
	detail.BillKey = getStringValue("bill_key")

	// Store specific fields
	detail.Store = getStringValue("store")
	detail.PaymentCode = getStringValue("payment_code")

	// E-wallet specific fields
	detail.Issuer = getStringValue("issuer")
	detail.Acquirer = getStringValue("acquirer")

	// Handle JSON arrays
	if vaNumbers, ok := notification["va_numbers"]; ok {
		jsonBytes, err := json.Marshal(vaNumbers)
		if err == nil {
			detail.VANumbers = model.JSON(jsonBytes)
		}
	}

	if paymentAmounts, ok := notification["payment_amounts"]; ok {
		jsonBytes, err := json.Marshal(paymentAmounts)
		if err == nil {
			detail.PaymentAmounts = model.JSON(jsonBytes)
		}
	}

	return detail
}

// Helper function to get string value from map with default
func getString(data map[string]interface{}, key string, defaultValue string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return defaultValue
}

func (s *subscriptionService) GetUserActiveSubscription(ctx *fiber.Ctx, userID uuid.UUID) (*model.UserSubscriptionResponse, error) {
	var subscription model.UserSubscription
	err := s.DB.WithContext(ctx.Context()).
		Joins("Plan").
		Where("user_id = ? AND end_date > ? AND is_active = ?", userID, time.Now(), true).
		First(&subscription).Error

	if err != nil {
		return nil, err
	}

	return s.toSubscriptionResponse(&subscription)
}

func (s *subscriptionService) toSubscriptionResponse(sub *model.UserSubscription) (*model.UserSubscriptionResponse, error) {
	var features map[string]bool
	// Initialize features map
	features = make(map[string]bool)

	// Only try to unmarshal if Features is not empty
	if sub.Plan.Features != "" {
		if err := json.Unmarshal([]byte(sub.Plan.Features), &features); err != nil {
			s.Log.Errorf("Failed to unmarshal features: %v", err)
			return nil, fmt.Errorf("invalid feature format: %w", err)
		}
	}

	return &model.UserSubscriptionResponse{
		ID:     sub.ID,
		UserID: sub.UserID,
		Plan: model.SubscriptionPlanResponse{
			ID:             sub.Plan.ID,
			Name:           sub.Plan.Name,
			Price:          sub.Plan.Price,
			PriceFormatted: formatCurrency(sub.Plan.Price),
			Features:       features,
			Description:    sub.Plan.Description,
			ValidityDays:   sub.Plan.ValidityDays,
			AIscanLimit:    sub.Plan.AIscanLimit,
		},
		AIscansUsed:   sub.AIscansUsed,
		StartDate:     sub.StartDate,
		EndDate:       sub.EndDate,
		IsActive:      sub.IsActive,
		PaymentMethod: sub.PaymentMethod,
		PaymentStatus: sub.PaymentStatus,
		CreatedAt:     sub.CreatedAt,
	}, nil
}

func (s *subscriptionService) CheckFeatureAccess(ctx *fiber.Ctx, userID uuid.UUID, feature string) (bool, error) {
	sub, err := s.GetUserActiveSubscription(ctx, userID)
	if err != nil {
		return false, nil
	}
	return sub.Plan.Features[feature], nil
}

func (s *subscriptionService) IncrementScanUsage(ctx *fiber.Ctx, userID uuid.UUID) error {
	return s.DB.WithContext(ctx.Context()).
		Model(&model.UserSubscription{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("ai_scans_used", gorm.Expr("ai_scans_used + 1")).
		Error
}

func (s *subscriptionService) GetRemainingScans(ctx *fiber.Ctx, userID uuid.UUID) (int, error) {
	sub, err := s.GetUserActiveSubscription(ctx, userID)
	if err != nil {
		return 0, err
	}
	remaining := sub.Plan.AIscanLimit - sub.AIscansUsed
	if remaining < 0 {
		return 0, nil
	}
	return remaining, nil
}
