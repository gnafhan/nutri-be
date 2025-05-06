package service

import (
	"app/src/config"
	"app/src/model"
	res "app/src/response"
	"app/src/utils"
	"app/src/validation"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TokenService interface {
	GenerateToken(c *fiber.Ctx, user *model.User, isProductTokenVerified bool, expires time.Time, tokenType string) (string, error)
	SaveToken(c *fiber.Ctx, token, userID, tokenType string, expires time.Time) error
	DeleteToken(c *fiber.Ctx, tokenType string, userID string) error
	DeleteAllToken(c *fiber.Ctx, userID string) error
	GetTokenByUserID(c *fiber.Ctx, tokenStr string) (*model.Token, error)
	GenerateAuthTokens(c *fiber.Ctx, user *model.User) (*res.Tokens, error)
	GenerateResetPasswordToken(c *fiber.Ctx, req *validation.ForgotPassword) (string, error)
	GenerateVerifyEmailToken(c *fiber.Ctx, user *model.User) (*string, error)
}

type tokenService struct {
	Log                 *logrus.Logger
	DB                  *gorm.DB
	Validate            *validator.Validate
	UserService         UserService
	SubscriptionService SubscriptionService
}

func NewTokenService(db *gorm.DB, validate *validator.Validate, userService UserService, subscriptionService SubscriptionService) TokenService {
	return &tokenService{
		Log:                 utils.Log,
		DB:                  db,
		Validate:            validate,
		UserService:         userService,
		SubscriptionService: subscriptionService,
	}
}

// ✅ Generate JWT dengan `userData`
func (s *tokenService) GenerateToken(c *fiber.Ctx, user *model.User, isProductTokenVerified bool, expires time.Time, tokenType string) (string, error) {
	// Get user's subscription features
	var subscriptionFeatures map[string]bool
	subscriptionFeatures = make(map[string]bool) // Initialize with empty map as default

	// Only try to get subscription if we have a context
	if c != nil {
		subscription, err := s.SubscriptionService.GetUserActiveSubscription(c, user.ID)
		if err == nil && subscription != nil {
			// User has an active subscription
			subscriptionFeatures = subscription.Plan.Features
		}
	}

	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"iat":  time.Now().Unix(),
		"exp":  expires.Unix(),
		"type": tokenType,
		"userData": map[string]interface{}{
			"id":                     user.ID,
			"name":                   user.Name,
			"email":                  user.Email,
			"role":                   user.Role,
			"verified_email":         user.VerifiedEmail,
			"birth_date":             user.BirthDate,
			"height":                 user.Height,
			"weight":                 user.Weight,
			"gender":                 user.Gender,
			"activity_level":         user.ActivityLevel,
			"medical_history":        user.MedicalHistory,
			"profile_picture":        user.ProfilePicture,
			"isProductTokenVerified": isProductTokenVerified,
			"subscriptionFeatures":   subscriptionFeatures,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

// ✅ Simpan Token ke Database
func (s *tokenService) SaveToken(c *fiber.Ctx, token, userID, tokenType string, expires time.Time) error {
	if err := s.DeleteToken(c, tokenType, userID); err != nil {
		return err
	}

	tokenDoc := &model.Token{
		Token:   token,
		UserID:  uuid.MustParse(userID),
		Type:    tokenType,
		Expires: expires,
	}

	return s.DB.WithContext(c.Context()).Create(tokenDoc).Error
}

// ✅ Hapus Token Berdasarkan Jenisnya
func (s *tokenService) DeleteToken(c *fiber.Ctx, tokenType string, userID string) error {
	return s.DB.WithContext(c.Context()).Where("type = ? AND user_id = ?", tokenType, userID).Delete(&model.Token{}).Error
}

// ✅ Hapus Semua Token User
func (s *tokenService) DeleteAllToken(c *fiber.Ctx, userID string) error {
	return s.DB.WithContext(c.Context()).Where("user_id = ?", userID).Delete(&model.Token{}).Error
}

// ✅ Ambil Token Berdasarkan User ID
func (s *tokenService) GetTokenByUserID(c *fiber.Ctx, tokenStr string) (*model.Token, error) {
	userID, err := utils.VerifyToken(tokenStr, config.JWTSecret, config.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	tokenDoc := new(model.Token)
	err = s.DB.WithContext(c.Context()).Where("token = ? AND user_id = ?", tokenStr, userID).First(tokenDoc).Error
	return tokenDoc, err
}

// ✅ Generate Access & Refresh Tokens
func (s *tokenService) GenerateAuthTokens(c *fiber.Ctx, user *model.User) (*res.Tokens, error) {
	// Cek apakah user memiliki Product Token yang aktif
	var isProductTokenVerified bool
	if err := s.DB.WithContext(c.Context()).Where("user_id = ?", user.ID).First(&model.ProductToken{}).Error; err == nil {
		isProductTokenVerified = true
	}

	// Generate Access Token
	accessTokenExpires := time.Now().UTC().Add(time.Minute * time.Duration(config.JWTAccessExp))
	accessToken, err := s.GenerateToken(c, user, isProductTokenVerified, accessTokenExpires, config.TokenTypeAccess)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshTokenExpires := time.Now().UTC().Add(time.Hour * 24 * time.Duration(config.JWTRefreshExp))
	refreshToken, err := s.GenerateToken(c, user, isProductTokenVerified, refreshTokenExpires, config.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	// Simpan Refresh Token ke Database
	if err = s.SaveToken(c, refreshToken, user.ID.String(), config.TokenTypeRefresh, refreshTokenExpires); err != nil {
		return nil, err
	}

	return &res.Tokens{
		Access: res.TokenExpires{
			Token:   accessToken,
			Expires: accessTokenExpires,
		},
		Refresh: res.TokenExpires{
			Token:   refreshToken,
			Expires: refreshTokenExpires,
		},
	}, nil
}

// ✅ Generate Token Reset Password
func (s *tokenService) GenerateResetPasswordToken(c *fiber.Ctx, req *validation.ForgotPassword) (string, error) {
	if err := s.Validate.Struct(req); err != nil {
		return "", err
	}

	user, err := s.UserService.GetUserByEmail(c, req.Email)
	if err != nil {
		return "", err
	}

	expires := time.Now().UTC().Add(time.Minute * time.Duration(config.JWTResetPasswordExp))
	return s.GenerateToken(c, user, false, expires, config.TokenTypeResetPassword)
}

// ✅ Generate Token Verifikasi Email
func (s *tokenService) GenerateVerifyEmailToken(c *fiber.Ctx, user *model.User) (*string, error) {
	expires := time.Now().UTC().Add(time.Minute * time.Duration(config.JWTVerifyEmailExp))
	verifyEmailToken, err := s.GenerateToken(c, user, false, expires, config.TokenTypeVerifyEmail)
	if err != nil {
		return nil, err
	}

	// Simpan token di database
	if err = s.SaveToken(c, verifyEmailToken, user.ID.String(), config.TokenTypeVerifyEmail, expires); err != nil {
		return nil, err
	}

	return &verifyEmailToken, nil
}
