package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductTokenService interface {
	GetProductTokenByUserID(c *fiber.Ctx, userId uuid.UUID) (*model.ProductToken, error)
	DeleteProductToken(c *fiber.Ctx, tokenID uuid.UUID) error
	VerifyProductToken(c *fiber.Ctx, query *validation.Token) error

	// Admin endpoints
	GetAllProductTokens(c *fiber.Ctx, query *validation.ProductTokenQuery) ([]model.ProductToken, error)
	CreateProductToken(c *fiber.Ctx, req *validation.CreateCustomToken) (*model.ProductToken, error)
	AdminDeleteProductToken(c *fiber.Ctx, tokenID uuid.UUID) error
	UpdateProductToken(c *fiber.Ctx, tokenID uuid.UUID, req *validation.UpdateProductToken) (*model.ProductToken, error)
}

type productTokenService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

// NewProductTokenService membuat instance service
func NewProductTokenService(db *gorm.DB, validate *validator.Validate) *productTokenService {
	return &productTokenService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *productTokenService) GetProductTokenByUserID(c *fiber.Ctx, userID uuid.UUID) (*model.ProductToken, error) {
	var productToken model.ProductToken
	err := s.DB.WithContext(c.Context()).
		Where("user_id = ?", userID).
		First(&productToken).Error

	if err != nil {
		return nil, err
	}
	return &productToken, nil
}

func (s *productTokenService) DeleteProductToken(c *fiber.Ctx, tokenID uuid.UUID) error {
	return s.DB.WithContext(c.Context()).
		Where("id = ?", tokenID).
		Delete(&model.ProductToken{}).Error
}

func (s *productTokenService) VerifyProductToken(c *fiber.Ctx, query *validation.Token) error {
	if err := s.Validate.Struct(query); err != nil {
		return err
	}

	userData := c.Locals("user")
	user, ok := userData.(*model.User)
	fmt.Println("User stored in Locals:", user)

	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}

	existingToken, _ := s.GetProductTokenByUserID(c, user.ID)
	if existingToken != nil {
		return fiber.NewError(fiber.StatusForbidden, "Can only be connected with 1 product token.")
	}

	var productToken model.ProductToken
	if err := s.DB.WithContext(c.Context()).
		Preload("SubscriptionPlan").
		Where("token = ? AND user_id IS NULL AND is_active = ?", query.Token, true).
		First(&productToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Invalid or already used product token")
		}
		return fiber.ErrInternalServerError
	}

	now := time.Now()

	updateData := map[string]interface{}{
		"user_id":      user.ID,
		"activated_at": &now,
	}

	if err := s.DB.WithContext(c.Context()).Model(&model.ProductToken{}).Where("id = ?", productToken.ID).Updates(updateData).Error; err != nil {
		s.Log.Errorf("Failed to update product token %s with user ID %s: %v", productToken.ID, user.ID, err)
		return fiber.ErrInternalServerError
	}

	productToken.UserID = user.ID
	productToken.ActivatedAt = &now

	// If product token has an associated subscription plan, create it for the user
	if productToken.SubscriptionPlanID != nil && productToken.SubscriptionPlan != nil {
		// Check if user already has an active subscription to this plan to avoid duplicates
		var existingUserSubscription model.UserSubscription
		err := s.DB.WithContext(c.Context()).
			Where("user_id = ? AND plan_id = ? AND is_active = ?", user.ID, productToken.SubscriptionPlanID, true).
			First(&existingUserSubscription).Error

		if err == nil {
			s.Log.Infof("User %s already has an active subscription to plan %s. Skipping creation.", user.ID, *productToken.SubscriptionPlanID)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// Check if user has a freemium subscription that needs to be upgraded
			var freemiumSubscription model.UserSubscription
			freemiumErr := s.DB.WithContext(c.Context()).
				Joins("JOIN subscription_plans ON user_subscriptions.plan_id = subscription_plans.id").
				Where("user_subscriptions.user_id = ? AND subscription_plans.name = ? AND user_subscriptions.is_active = ?", user.ID, "Freemium Trial", true).
				First(&freemiumSubscription).Error

			if freemiumErr == nil {
				// User has an active freemium subscription, deactivate it
				s.Log.Infof("Deactivating freemium subscription for user %s to upgrade to product token plan", user.ID)
				if err := s.DB.WithContext(c.Context()).
					Model(&freemiumSubscription).
					Update("is_active", false).Error; err != nil {
					s.Log.Errorf("Failed to deactivate freemium subscription for user %s: %v", user.ID, err)
				}
			}

			// Create the new subscription from product token
			userSubscription := model.UserSubscription{
				UserID:        user.ID,
				PlanID:        *productToken.SubscriptionPlanID,
				StartDate:     time.Now(),
				EndDate:       time.Now().AddDate(0, 0, productToken.SubscriptionPlan.ValidityDays),
				PaymentMethod: "product_token",
				PaymentStatus: "success", // Assuming token verification implies successful "payment"
				IsActive:      true,
				TransactionID: fmt.Sprintf("TOKEN-%s", productToken.ID.String()), // Link to product token
			}
			if err := s.DB.WithContext(c.Context()).Create(&userSubscription).Error; err != nil {
				s.Log.Errorf("Failed to create subscription for user %s with plan %s: %v", user.ID, *productToken.SubscriptionPlanID, err)
				// Decide if this should be a hard error or just a warning. For now, log and continue.
			} else {
				s.Log.Infof("Successfully created product token subscription for user %s with plan %s", user.ID, *productToken.SubscriptionPlanID)
			}
		} else {
			// Some other DB error occurred when checking for existing subscription
			s.Log.Errorf("Error checking for existing subscription for user %s, plan %s: %v", user.ID, *productToken.SubscriptionPlanID, err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product token verified successfully",
	})
}

// Admin functions

func (s *productTokenService) GetAllProductTokens(c *fiber.Ctx, query *validation.ProductTokenQuery) ([]model.ProductToken, error) {
	var tokens []model.ProductToken
	db := s.DB.WithContext(c.Context()).Preload("SubscriptionPlan")

	if query != nil && query.WithUser {
		db = db.Preload("User").Preload("CreatedBy")
	}

	if err := db.Find(&tokens).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *productTokenService) CreateProductToken(c *fiber.Ctx, req *validation.CreateCustomToken) (*model.ProductToken, error) {
	userData := c.Locals("user")
	admin, ok := userData.(*model.User)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}

	if err := s.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request data")
	}

	// Jika token kustom disediakan, gunakan itu; jika tidak, hasilkan token acak
	token := req.Token
	if token == "" {
		token = utils.GenerateRandomString(16)
	}

	// Periksa apakah token sudah ada
	var existingCount int64
	if err := s.DB.WithContext(c.Context()).Model(&model.ProductToken{}).Where("token = ?", token).Count(&existingCount).Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if existingCount > 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Token already exists")
	}

	productToken := model.ProductToken{
		Token:       token,
		CreatedByID: admin.ID,
		IsActive:    req.IsActive,
	}

	if req.SubscriptionPlanID != nil && *req.SubscriptionPlanID != "" {
		planID, err := uuid.Parse(*req.SubscriptionPlanID)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid SubscriptionPlanID format")
		}
		// Verify if the plan ID exists
		var plan model.SubscriptionPlan
		if err := s.DB.WithContext(c.Context()).First(&plan, "id = ?", planID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fiber.NewError(fiber.StatusNotFound, "Subscription plan not found")
			}
			return nil, fiber.ErrInternalServerError
		}
		productToken.SubscriptionPlanID = &planID
	}

	if err := s.DB.WithContext(c.Context()).Create(&productToken).Error; err != nil {
		return nil, err
	}

	// Reload the token with the creator information
	if err := s.DB.WithContext(c.Context()).Preload("CreatedBy").Preload("SubscriptionPlan").First(&productToken, productToken.ID).Error; err != nil {
		s.Log.Warnf("Unable to load creator or subscription plan information: %v", err)
	}

	return &productToken, nil
}

func (s *productTokenService) AdminDeleteProductToken(c *fiber.Ctx, tokenID uuid.UUID) error {
	return s.DB.WithContext(c.Context()).
		Where("id = ?", tokenID).
		Delete(&model.ProductToken{}).Error
}

func (s *productTokenService) UpdateProductToken(c *fiber.Ctx, tokenID uuid.UUID, req *validation.UpdateProductToken) (*model.ProductToken, error) {
	var productToken model.ProductToken

	if err := s.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request data: "+err.Error())
	}

	if err := s.DB.WithContext(c.Context()).First(&productToken, "id = ?", tokenID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Product token not found")
		}
		return nil, fiber.ErrInternalServerError
	}

	// Update token string if provided and different
	if req.Token != nil && *req.Token != "" && *req.Token != productToken.Token {
		// Check if the new token already exists for another record
		var existingTokenCount int64
		if err := s.DB.WithContext(c.Context()).Model(&model.ProductToken{}).Where("token = ? AND id <> ?", *req.Token, tokenID).Count(&existingTokenCount).Error; err != nil {
			return nil, fiber.ErrInternalServerError
		}
		if existingTokenCount > 0 {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Token already exists")
		}
		productToken.Token = *req.Token
	}

	// Update IsActive if provided
	if req.IsActive != nil {
		productToken.IsActive = *req.IsActive
	}

	// Update SubscriptionPlanID if provided
	if req.SubscriptionPlanID != nil {
		if *req.SubscriptionPlanID == "" { // Admin wants to remove the plan
			productToken.SubscriptionPlanID = nil
		} else {
			planID, err := uuid.Parse(*req.SubscriptionPlanID)
			if err != nil {
				return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid SubscriptionPlanID format")
			}
			// Verify if the plan ID exists
			var plan model.SubscriptionPlan
			if err := s.DB.WithContext(c.Context()).First(&plan, "id = ?", planID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, fiber.NewError(fiber.StatusNotFound, "Subscription plan not found")
				}
				return nil, fiber.ErrInternalServerError
			}
			productToken.SubscriptionPlanID = &planID
		}
	}

	if err := s.DB.WithContext(c.Context()).Save(&productToken).Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	// Reload the token with potentially updated relations
	if err := s.DB.WithContext(c.Context()).Preload("CreatedBy").Preload("User").Preload("SubscriptionPlan").First(&productToken, productToken.ID).Error; err != nil {
		s.Log.Warnf("Unable to load full product token information after update: %v", err)
		// Not returning error here, as the main update succeeded
	}

	return &productToken, nil
}
