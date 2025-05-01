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
}

type PaymentResponse struct {
	TransactionID string
	Status        string
}

type SubscriptionService interface {
	GetAllPlans(ctx *fiber.Ctx) ([]model.SubscriptionPlanResponse, error)
	PurchasePlan(ctx *fiber.Ctx, userID uuid.UUID, planID uuid.UUID, paymentMethod string) (*model.UserSubscriptionResponse, error)
	GetUserActiveSubscription(ctx *fiber.Ctx, userID uuid.UUID) (*model.UserSubscriptionResponse, error)
	CheckFeatureAccess(ctx *fiber.Ctx, userID uuid.UUID, feature string) (bool, error)
	IncrementScanUsage(ctx *fiber.Ctx, userID uuid.UUID) error
	GetRemainingScans(ctx *fiber.Ctx, userID uuid.UUID) (int, error)
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
			IsRecommended:  plan.Name == "Sehat",
			Description:    plan.Description,
			ValidityDays:   plan.ValidityDays,
			AIscanLimit:    plan.AIscanLimit,
		})
	}

	return responses, nil
}

func (s *subscriptionService) PurchasePlan(ctx *fiber.Ctx, userID uuid.UUID, planID uuid.UUID, paymentMethod string) (*model.UserSubscriptionResponse, error) {
	var plan model.SubscriptionPlan
	if err := s.DB.WithContext(ctx.Context()).First(&plan, "id = ?", planID).Error; err != nil {
		return nil, errors.New("subscription plan not found")
	}

	paymentResp, err := s.Payment.Charge(plan.Price, paymentMethod)
	if err != nil {
		return nil, fmt.Errorf("payment failed: %w", err)
	}

	subscription := model.UserSubscription{
		UserID:        userID,
		PlanID:        planID,
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 0, plan.ValidityDays),
		PaymentMethod: paymentMethod,
		TransactionID: paymentResp.TransactionID,
	}

	err = s.DB.WithContext(ctx.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.UserSubscription{}).
			Where("user_id = ? AND is_active = ?", userID, true).
			Update("is_active", false).Error; err != nil {
			return err
		}
		return tx.Create(&subscription).Error
	})

	if err != nil {
		if refundErr := s.Payment.Refund(paymentResp.TransactionID); refundErr != nil {
			s.Log.Errorf("Failed to refund payment: %v", refundErr)
		}
		return nil, err
	}

	return s.toSubscriptionResponse(&subscription)
}

func (s *subscriptionService) GetUserActiveSubscription(ctx *fiber.Ctx, userID uuid.UUID) (*model.UserSubscriptionResponse, error) {
	var subscription model.UserSubscription
	err := s.DB.WithContext(ctx.Context()).
		Joins("Plan").
		Where("user_id = ? AND is_active = ? AND end_date > ?", userID, true, time.Now()).
		First(&subscription).Error

	if err != nil {
		return nil, err
	}

	return s.toSubscriptionResponse(&subscription)
}

func (s *subscriptionService) toSubscriptionResponse(sub *model.UserSubscription) (*model.UserSubscriptionResponse, error) {
	var features map[string]bool
	if err := json.Unmarshal([]byte(sub.Plan.Features), &features); err != nil {
		return nil, err
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
