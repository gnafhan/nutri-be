package service

import (
	"app/src/model"
	"app/src/validation"
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
	CreateFreemiumSubscription(ctx *fiber.Ctx, userID uuid.UUID) error

	// Admin-related methods
	GetAllUserSubscriptions(ctx *fiber.Ctx, query *validation.SubscriptionQuery) ([]model.UserSubscriptionResponse, int64, error)
	GetUserSubscriptionByID(ctx *fiber.Ctx, subscriptionID uuid.UUID) (*model.UserSubscriptionResponse, error)
	GetAllSubscriptionPlansWithUsers(ctx *fiber.Ctx, withUsers bool) ([]model.SubscriptionPlanWithUsers, error)
	UpdateUserSubscription(ctx *fiber.Ctx, subscriptionID uuid.UUID, req *validation.UpdateSubscription) (*model.UserSubscriptionResponse, error)
	DeleteUserSubscription(ctx *fiber.Ctx, subscriptionID uuid.UUID) error
	GetTransactionsBySubscriptionID(ctx *fiber.Ctx, subscriptionID uuid.UUID) ([]model.TransactionDetail, error)
	UpdatePaymentStatus(ctx *fiber.Ctx, subscriptionID uuid.UUID, status string) (*model.UserSubscriptionResponse, error)
	GetAllTransactions(ctx *fiber.Ctx, page, limit int) ([]model.TransactionDetail, int64, error)
	GetTransactionByID(ctx *fiber.Ctx, transactionID uuid.UUID) (*model.TransactionDetail, error)
	GetSubscriptionPlanByID(ctx *fiber.Ctx, planID uuid.UUID) (*model.SubscriptionPlan, error)
	UpdateSubscriptionPlan(ctx *fiber.Ctx, planID uuid.UUID, req *validation.UpdateSubscriptionPlan) (*model.SubscriptionPlan, error)
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
		Preload("Plan").
		Where("user_subscriptions.user_id = ? AND user_subscriptions.end_date > ? AND user_subscriptions.is_active = ? AND user_subscriptions.payment_status = ?", userID, time.Now(), true, "success").
		First(&subscription).Error

	if err != nil {
		return nil, err
	}

	return s.toSubscriptionResponse(&subscription)
}

func (s *subscriptionService) toSubscriptionResponse(sub *model.UserSubscription) (*model.UserSubscriptionResponse, error) {
	// Initialize features map
	features := make(map[string]bool)

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

// GetAllUserSubscriptions retrieves all user subscriptions with pagination and filtering
func (s *subscriptionService) GetAllUserSubscriptions(ctx *fiber.Ctx, query *validation.SubscriptionQuery) ([]model.UserSubscriptionResponse, int64, error) {
	var subscriptions []model.UserSubscription
	var totalResults int64

	db := s.DB.WithContext(ctx.Context()).
		Preload("Plan").
		Preload("User")

	// Apply status filter if provided
	if query.Status != "" {
		db = db.Where("user_subscriptions.payment_status = ?", query.Status)
	}

	// Count total results
	if err := db.Model(&model.UserSubscription{}).Count(&totalResults).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if err := db.
		Offset((query.Page - 1) * query.Limit).
		Limit(query.Limit).
		Order("user_subscriptions.created_at DESC").
		Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}

	// Convert to response format
	var responses []model.UserSubscriptionResponse
	for _, sub := range subscriptions {
		response, err := s.toSubscriptionResponse(&sub)
		if err != nil {
			return nil, 0, err
		}
		responses = append(responses, *response)
	}

	return responses, totalResults, nil
}

// GetUserSubscriptionByID retrieves a specific user subscription by ID
func (s *subscriptionService) GetUserSubscriptionByID(ctx *fiber.Ctx, subscriptionID uuid.UUID) (*model.UserSubscriptionResponse, error) {
	var subscription model.UserSubscription

	if err := s.DB.WithContext(ctx.Context()).
		Preload("Plan").
		Preload("User").
		Where("user_subscriptions.id = ?", subscriptionID).
		First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Subscription not found")
		}
		return nil, err
	}

	return s.toSubscriptionResponse(&subscription)
}

// GetAllSubscriptionPlansWithUsers retrieves all subscription plans with their users
func (s *subscriptionService) GetAllSubscriptionPlansWithUsers(ctx *fiber.Ctx, withUsers bool) ([]model.SubscriptionPlanWithUsers, error) {
	var plans []model.SubscriptionPlan

	if err := s.DB.WithContext(ctx.Context()).Find(&plans).Error; err != nil {
		return nil, err
	}

	var result []model.SubscriptionPlanWithUsers

	for _, plan := range plans {
		// Parse features
		var features map[string]bool
		if err := json.Unmarshal([]byte(plan.Features), &features); err != nil {
			return nil, err
		}

		planWithUsers := model.SubscriptionPlanWithUsers{
			ID:             plan.ID,
			Name:           plan.Name,
			Price:          plan.Price,
			PriceFormatted: formatCurrency(plan.Price),
			Description:    plan.Description,
			AIscanLimit:    plan.AIscanLimit,
			ValidityDays:   plan.ValidityDays,
			Features:       features,
			IsActive:       plan.IsActive,
		}

		// Count users for this plan
		var userCount int64
		if err := s.DB.WithContext(ctx.Context()).
			Model(&model.UserSubscription{}).
			Where("plan_id = ? AND is_active = ?", plan.ID, true).
			Count(&userCount).Error; err != nil {
			return nil, err
		}

		planWithUsers.UserCount = int(userCount)

		// Get users if requested
		if withUsers {
			var subscriptions []model.UserSubscription
			if err := s.DB.WithContext(ctx.Context()).
				Preload("User").
				Where("plan_id = ? AND is_active = ?", plan.ID, true).
				Find(&subscriptions).Error; err != nil {
				return nil, err
			}

			for _, sub := range subscriptions {
				response, err := s.toSubscriptionResponse(&sub)
				if err != nil {
					return nil, err
				}
				planWithUsers.Users = append(planWithUsers.Users, *response)
			}
		}

		result = append(result, planWithUsers)
	}

	return result, nil
}

// UpdateUserSubscription updates a user subscription
func (s *subscriptionService) UpdateUserSubscription(ctx *fiber.Ctx, subscriptionID uuid.UUID, req *validation.UpdateSubscription) (*model.UserSubscriptionResponse, error) {
	var subscription model.UserSubscription

	if err := s.DB.WithContext(ctx.Context()).
		Preload("Plan").
		Where("user_subscriptions.id = ?", subscriptionID).
		First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Subscription not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.PlanID != nil {
		// Verify the plan exists
		var plan model.SubscriptionPlan
		if err := s.DB.WithContext(ctx.Context()).First(&plan, "id = ?", *req.PlanID).Error; err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid plan ID")
		}
		subscription.PlanID = *req.PlanID
	}

	if req.IsActive != nil {
		subscription.IsActive = *req.IsActive
	}

	if req.AIscansUsed != nil {
		subscription.AIscansUsed = *req.AIscansUsed
	}

	if req.StartDate != nil {
		subscription.StartDate = *req.StartDate
	}

	if req.EndDate != nil {
		subscription.EndDate = *req.EndDate
	}

	if req.PaymentMethod != nil {
		subscription.PaymentMethod = *req.PaymentMethod
	}

	// Save changes
	if err := s.DB.WithContext(ctx.Context()).Save(&subscription).Error; err != nil {
		return nil, err
	}

	// Refresh subscription data
	if err := s.DB.WithContext(ctx.Context()).
		Preload("Plan").
		Where("user_subscriptions.id = ?", subscriptionID).
		First(&subscription).Error; err != nil {
		return nil, err
	}

	return s.toSubscriptionResponse(&subscription)
}

// DeleteUserSubscription deletes a user subscription
func (s *subscriptionService) DeleteUserSubscription(ctx *fiber.Ctx, subscriptionID uuid.UUID) error {
	var subscription model.UserSubscription

	if err := s.DB.WithContext(ctx.Context()).
		Where("id = ?", subscriptionID).
		First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Subscription not found")
		}
		return err
	}

	// Delete subscription
	if err := s.DB.WithContext(ctx.Context()).Delete(&subscription).Error; err != nil {
		return err
	}

	return nil
}

// GetTransactionsBySubscriptionID retrieves all transactions for a subscription
func (s *subscriptionService) GetTransactionsBySubscriptionID(ctx *fiber.Ctx, subscriptionID uuid.UUID) ([]model.TransactionDetail, error) {
	var transactions []model.TransactionDetail

	if err := s.DB.WithContext(ctx.Context()).
		Where("user_subscription_id = ?", subscriptionID).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

// UpdatePaymentStatus updates the payment status of a subscription
func (s *subscriptionService) UpdatePaymentStatus(ctx *fiber.Ctx, subscriptionID uuid.UUID, status string) (*model.UserSubscriptionResponse, error) {
	var subscription model.UserSubscription

	if err := s.DB.WithContext(ctx.Context()).
		Joins("Plan").
		Where("user_subscriptions.id = ?", subscriptionID).
		First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Subscription not found")
		}
		return nil, err
	}

	// Update payment status
	subscription.PaymentStatus = status

	// If status is success, activate the subscription
	if status == "success" {
		subscription.IsActive = true
	} else if status == "failed" {
		subscription.IsActive = false
	}

	// Save changes
	if err := s.DB.WithContext(ctx.Context()).Save(&subscription).Error; err != nil {
		return nil, err
	}

	// Create transaction record
	transactionDetail := &model.TransactionDetail{
		UserSubscriptionID: subscription.ID,
		OrderID:            subscription.TransactionID,
		TransactionStatus:  status,
		TransactionTime:    time.Now(),
		GrossAmount:        fmt.Sprintf("%d", subscription.Plan.Price),
		Currency:           "IDR",
	}

	if err := s.DB.WithContext(ctx.Context()).Create(&transactionDetail).Error; err != nil {
		s.Log.Warnf("Failed to create transaction record: %v", err)
		// Continue even if transaction record creation fails
	}

	return s.toSubscriptionResponse(&subscription)
}

// GetAllTransactions retrieves all transaction logs with pagination
func (s *subscriptionService) GetAllTransactions(ctx *fiber.Ctx, page, limit int) ([]model.TransactionDetail, int64, error) {
	var transactions []model.TransactionDetail
	var totalResults int64

	// Count total results
	if err := s.DB.WithContext(ctx.Context()).
		Model(&model.TransactionDetail{}).
		Count(&totalResults).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if err := s.DB.WithContext(ctx.Context()).
		Preload("UserSubscription").
		Preload("UserSubscription.User").
		Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, totalResults, nil
}

func (s *subscriptionService) GetTransactionByID(ctx *fiber.Ctx, transactionID uuid.UUID) (*model.TransactionDetail, error) {
	var transaction model.TransactionDetail

	if err := s.DB.WithContext(ctx.Context()).
		Preload("UserSubscription").
		Preload("UserSubscription.User").
		Preload("UserSubscription.Plan").
		Where("id = ?", transactionID).
		First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Transaction not found")
		}
		return nil, err
	}

	return &transaction, nil
}

func (s *subscriptionService) GetSubscriptionPlanByID(ctx *fiber.Ctx, planID uuid.UUID) (*model.SubscriptionPlan, error) {
	var plan model.SubscriptionPlan

	if err := s.DB.WithContext(ctx.Context()).First(&plan, "id = ?", planID).Error; err != nil {
		return nil, err
	}

	return &plan, nil
}

func (s *subscriptionService) UpdateSubscriptionPlan(ctx *fiber.Ctx, planID uuid.UUID, req *validation.UpdateSubscriptionPlan) (*model.SubscriptionPlan, error) {
	var plan model.SubscriptionPlan

	if err := s.DB.WithContext(ctx.Context()).First(&plan, "id = ?", planID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Subscription plan not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		plan.Name = *req.Name
	}

	if req.Price != nil {
		plan.Price = *req.Price
	}

	if req.ValidityDays != nil {
		plan.ValidityDays = *req.ValidityDays
	}

	if req.AIscanLimit != nil {
		plan.AIscanLimit = *req.AIscanLimit
	}

	if req.IsActive != nil {
		plan.IsActive = *req.IsActive
	}

	if req.Description != nil {
		plan.Description = *req.Description
	}

	// Update features if provided
	if req.Features != nil {
		featuresJSON, err := json.Marshal(req.Features)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid features format")
		}
		plan.Features = string(featuresJSON)
	}

	// Save changes
	if err := s.DB.WithContext(ctx.Context()).Save(&plan).Error; err != nil {
		return nil, err
	}

	return &plan, nil
}

func (s *subscriptionService) CreateFreemiumSubscription(ctx *fiber.Ctx, userID uuid.UUID) error {
	// Check if user already has an active subscription
	existingSubscription, err := s.GetUserActiveSubscription(ctx, userID)
	if err == nil && existingSubscription != nil {
		s.Log.Infof("User %s already has an active subscription, skipping freemium creation", userID.String())
		return nil
	}

	// Find the Freemium Trial plan
	var freemiumPlan model.SubscriptionPlan
	if err := s.DB.WithContext(ctx.Context()).First(&freemiumPlan, "name = ? AND is_active = ?", "Freemium Trial", true).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "Freemium Trial plan not found")
		}
		return err
	}

	// Create the freemium subscription
	now := time.Now()
	endDate := now.AddDate(0, 0, 14) // 14 days from now

	freemiumSubscription := &model.UserSubscription{
		UserID:        userID,
		PlanID:        freemiumPlan.ID,
		StartDate:     now,
		EndDate:       endDate,
		IsActive:      true,
		PaymentMethod: "freemium_trial",
		PaymentStatus: "completed",
		AIscansUsed:   0,
	}

	if err := s.DB.WithContext(ctx.Context()).Create(freemiumSubscription).Error; err != nil {
		s.Log.Errorf("Failed to create freemium subscription for user %s: %v", userID.String(), err)
		return err
	}

	s.Log.Infof("Successfully created freemium subscription for user %s, expires at %s", userID.String(), endDate.Format(time.RFC3339))
	return nil
}
