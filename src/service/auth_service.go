package service

import (
	"app/src/config"
	"app/src/model"
	"app/src/response"
	"app/src/utils"
	"app/src/validation"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(c *fiber.Ctx, req *validation.Register) (*model.User, error)
	Login(c *fiber.Ctx, req *validation.Login) (*model.User, error)
	Logout(c *fiber.Ctx, req *validation.Logout) error
	RefreshAuth(c *fiber.Ctx, req *validation.RefreshToken) (*response.Tokens, error)
	ResetPassword(c *fiber.Ctx, query *validation.Token, req *validation.UpdatePassOrVerify) error
	VerifyEmail(c *fiber.Ctx, query *validation.Token) error
}

type authService struct {
	Log                 *logrus.Logger
	DB                  *gorm.DB
	Validate            *validator.Validate
	UserService         UserService
	TokenService        TokenService
	SubscriptionService SubscriptionService
}

func NewAuthService(
	db *gorm.DB, validate *validator.Validate, userService UserService, tokenService TokenService, subscriptionService SubscriptionService,
) AuthService {
	return &authService{
		Log:                 utils.Log,
		DB:                  db,
		Validate:            validate,
		UserService:         userService,
		TokenService:        tokenService,
		SubscriptionService: subscriptionService,
	}
}

func (s *authService) Register(c *fiber.Ctx, req *validation.Register) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.Log.Errorf("Failed to hash password: %+v", err)
		return nil, err
	}

	user := &model.User{
		Name:           req.Name,
		Email:          req.Email,
		Password:       hashedPassword,
		BirthDate:      &req.BirthDate,
		Height:         &req.Height,
		Weight:         &req.Weight,
		Gender:         &req.Gender,
		ActivityLevel:  &req.ActivityLevel,
		MedicalHistory: &req.MedicalHistory,
	}

	// Mulai transaksi database
	tx := s.DB.WithContext(c.Context()).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Simpan user
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fiber.NewError(fiber.StatusConflict, "Email already taken")
		}
		s.Log.Errorf("Failed to create user: %+v", err)
		return nil, err
	}

	userWeightHeight := &model.UsersWeightHeightHistory{
		ID:         uuid.New(),
		UserID:     user.ID,
		Weight:     req.Weight,
		Height:     req.Height,
		RecordedAt: time.Now(),
	}

	// Simpan userWeightHeight
	if err := tx.Create(userWeightHeight).Error; err != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to add weight height to database: %+v", err)
		return nil, err
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(c *fiber.Ctx, req *validation.Login) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	user, err := s.UserService.GetUserByEmail(c, req.Email)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	return user, nil
}

func (s *authService) Logout(c *fiber.Ctx, req *validation.Logout) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	token, err := s.TokenService.GetTokenByUserID(c, req.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Token not found")
	}

	err = s.TokenService.DeleteToken(c, config.TokenTypeRefresh, token.UserID.String())

	return err
}

func (s *authService) RefreshAuth(c *fiber.Ctx, req *validation.RefreshToken) (*response.Tokens, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	token, err := s.TokenService.GetTokenByUserID(c, req.RefreshToken)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	user, err := s.UserService.GetUserByID(c, token.UserID.String())
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	newTokens, err := s.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return newTokens, err
}

func (s *authService) ResetPassword(c *fiber.Ctx, query *validation.Token, req *validation.UpdatePassOrVerify) error {
	if err := s.Validate.Struct(query); err != nil {
		return err
	}

	userID, err := utils.VerifyToken(query.Token, config.JWTSecret, config.TokenTypeResetPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Token")
	}

	user, err := s.UserService.GetUserByID(c, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Password reset failed")
	}

	if errUpdate := s.UserService.UpdatePassOrVerify(c, req, user.ID.String()); errUpdate != nil {
		return errUpdate
	}

	if errToken := s.TokenService.DeleteToken(c, config.TokenTypeResetPassword, user.ID.String()); errToken != nil {
		return errToken
	}

	return nil
}

func (s *authService) VerifyEmail(c *fiber.Ctx, query *validation.Token) error {
	if err := s.Validate.Struct(query); err != nil {
		return err
	}

	userID, err := utils.VerifyToken(query.Token, config.JWTSecret, config.TokenTypeVerifyEmail)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Token")
	}

	user, err := s.UserService.GetUserByID(c, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Verify email failed")
	}

	// Start transaction for email verification and freemium subscription creation
	tx := s.DB.WithContext(c.Context()).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete verification token
	if errToken := s.TokenService.DeleteToken(c, config.TokenTypeVerifyEmail, user.ID.String()); errToken != nil {
		tx.Rollback()
		return errToken
	}

	// Update user email verification status
	updateBody := &validation.UpdatePassOrVerify{
		VerifiedEmail: true,
	}

	if errUpdate := s.UserService.UpdatePassOrVerify(c, updateBody, user.ID.String()); errUpdate != nil {
		tx.Rollback()
		return errUpdate
	}

	// Commit transaction first
	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("Failed to commit email verification transaction: %v", err)
		return err
	}

	// Create freemium subscription for new users (after successful email verification)
	// This is done outside the transaction to avoid dependency issues
	if errFreemium := s.SubscriptionService.CreateFreemiumSubscription(c, user.ID); errFreemium != nil {
		// Log the error but don't fail the email verification
		s.Log.Errorf("Failed to create freemium subscription for user %s: %v", user.ID.String(), errFreemium)
	}

	return nil
}
