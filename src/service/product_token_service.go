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
		Where("token = ? AND user_id IS NULL", query.Token).
		First(&productToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Invalid or already used product token")
		}
		return fiber.ErrInternalServerError
	}

	now := time.Now()
	productToken.UserID = user.ID
	productToken.ActivatedAt = &now

	if err := s.DB.WithContext(c.Context()).Save(&productToken).Error; err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product token verified successfully",
	})
}
