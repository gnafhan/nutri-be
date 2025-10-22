package middleware_test

import (
	"app/src/config"
	"app/src/middleware"
	"app/src/model"
	"app/src/response"
	"app/src/validation"
	"app/test/fixture"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserService is a mock implementation of UserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(c *fiber.Ctx, id string) (*model.User, error) {
	args := m.Called(c, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetUsers(c *fiber.Ctx, params *validation.QueryUser) ([]model.User, int64, error) {
	args := m.Called(c, params)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) GetUserByEmail(c *fiber.Ctx, email string) (*model.User, error) {
	args := m.Called(c, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) CreateUser(c *fiber.Ctx, req *validation.CreateUser) (*model.User, error) {
	args := m.Called(c, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdatePassOrVerify(c *fiber.Ctx, req *validation.UpdatePassOrVerify, id string) error {
	args := m.Called(c, req, id)
	return args.Error(0)
}

func (m *MockUserService) UpdateUser(c *fiber.Ctx, req *validation.UpdateUser, id string) (*model.User, error) {
	args := m.Called(c, req, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(c *fiber.Ctx, id string) error {
	args := m.Called(c, id)
	return args.Error(0)
}

func (m *MockUserService) CreateGoogleUser(c *fiber.Ctx, req *validation.GoogleLogin) (*model.User, error) {
	args := m.Called(c, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetUserStatistics(c *fiber.Ctx, userID string) (*response.UserStatistics, error) {
	args := m.Called(c, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.UserStatistics), args.Error(1)
}

// MockProductTokenService is a mock implementation of ProductTokenService interface
type MockProductTokenService struct {
	mock.Mock
}

func (m *MockProductTokenService) GetProductTokenByUserID(c *fiber.Ctx, userID uuid.UUID) (*model.ProductToken, error) {
	args := m.Called(c, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductToken), args.Error(1)
}

func (m *MockProductTokenService) DeleteProductToken(c *fiber.Ctx, tokenID uuid.UUID) error {
	args := m.Called(c, tokenID)
	return args.Error(0)
}

func (m *MockProductTokenService) VerifyProductToken(c *fiber.Ctx, query *validation.Token) error {
	args := m.Called(c, query)
	return args.Error(0)
}

func (m *MockProductTokenService) GetAllProductTokens(c *fiber.Ctx, query *validation.ProductTokenQuery) ([]model.ProductToken, error) {
	args := m.Called(c, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ProductToken), args.Error(1)
}

func (m *MockProductTokenService) CreateProductToken(c *fiber.Ctx, req *validation.CreateCustomToken) (*model.ProductToken, error) {
	args := m.Called(c, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductToken), args.Error(1)
}

func (m *MockProductTokenService) AdminDeleteProductToken(c *fiber.Ctx, tokenID uuid.UUID) error {
	args := m.Called(c, tokenID)
	return args.Error(0)
}

func (m *MockProductTokenService) UpdateProductToken(c *fiber.Ctx, tokenID uuid.UUID, req *validation.UpdateProductToken) (*model.ProductToken, error) {
	args := m.Called(c, tokenID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductToken), args.Error(1)
}

// MockSubscriptionService is a mock implementation of SubscriptionService interface
type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) GetUserActiveSubscription(c *fiber.Ctx, userID uuid.UUID) (*model.UserSubscriptionResponse, error) {
	args := m.Called(c, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserSubscriptionResponse), args.Error(1)
}

func (m *MockSubscriptionService) CheckFeatureAccess(c *fiber.Ctx, userID uuid.UUID, feature string) (bool, error) {
	args := m.Called(c, userID, feature)
	return args.Bool(0), args.Error(1)
}

func (m *MockSubscriptionService) CreateFreemiumSubscription(c *fiber.Ctx, userID uuid.UUID) error {
	args := m.Called(c, userID)
	return args.Error(0)
}

func (m *MockSubscriptionService) GetAllPlans(c *fiber.Ctx) ([]model.SubscriptionPlanResponse, error) {
	args := m.Called(c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.SubscriptionPlanResponse), args.Error(1)
}

func (m *MockSubscriptionService) PurchasePlan(c *fiber.Ctx, userID uuid.UUID, planID uuid.UUID, paymentMethod string) (*model.PaymentResponse, error) {
	args := m.Called(c, userID, planID, paymentMethod)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PaymentResponse), args.Error(1)
}

func (m *MockSubscriptionService) IncrementScanUsage(c *fiber.Ctx, userID uuid.UUID) error {
	args := m.Called(c, userID)
	return args.Error(0)
}

func (m *MockSubscriptionService) GetRemainingScans(c *fiber.Ctx, userID uuid.UUID) (int, error) {
	args := m.Called(c, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockSubscriptionService) HandlePaymentNotification(c *fiber.Ctx, notificationData []byte) error {
	args := m.Called(c, notificationData)
	return args.Error(0)
}

// Admin methods (stubbed for completeness)
func (m *MockSubscriptionService) GetAllUserSubscriptions(c *fiber.Ctx, query *validation.SubscriptionQuery) ([]model.UserSubscriptionResponse, int64, error) {
	args := m.Called(c, query)
	return args.Get(0).([]model.UserSubscriptionResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockSubscriptionService) GetUserSubscriptionByID(c *fiber.Ctx, subscriptionID uuid.UUID) (*model.UserSubscriptionResponse, error) {
	args := m.Called(c, subscriptionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserSubscriptionResponse), args.Error(1)
}

func (m *MockSubscriptionService) GetAllSubscriptionPlansWithUsers(c *fiber.Ctx, withUsers bool) ([]model.SubscriptionPlanWithUsers, error) {
	args := m.Called(c, withUsers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.SubscriptionPlanWithUsers), args.Error(1)
}

func (m *MockSubscriptionService) UpdateUserSubscription(c *fiber.Ctx, subscriptionID uuid.UUID, req *validation.UpdateSubscription) (*model.UserSubscriptionResponse, error) {
	args := m.Called(c, subscriptionID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserSubscriptionResponse), args.Error(1)
}

func (m *MockSubscriptionService) DeleteUserSubscription(c *fiber.Ctx, subscriptionID uuid.UUID) error {
	args := m.Called(c, subscriptionID)
	return args.Error(0)
}

func (m *MockSubscriptionService) GetTransactionsBySubscriptionID(c *fiber.Ctx, subscriptionID uuid.UUID) ([]model.TransactionDetail, error) {
	args := m.Called(c, subscriptionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TransactionDetail), args.Error(1)
}

func (m *MockSubscriptionService) UpdatePaymentStatus(c *fiber.Ctx, subscriptionID uuid.UUID, status string) (*model.UserSubscriptionResponse, error) {
	args := m.Called(c, subscriptionID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserSubscriptionResponse), args.Error(1)
}

func (m *MockSubscriptionService) GetAllTransactions(c *fiber.Ctx, page, limit int) ([]model.TransactionDetail, int64, error) {
	args := m.Called(c, page, limit)
	return args.Get(0).([]model.TransactionDetail), args.Get(1).(int64), args.Error(2)
}

func (m *MockSubscriptionService) GetTransactionByID(c *fiber.Ctx, transactionID uuid.UUID) (*model.TransactionDetail, error) {
	args := m.Called(c, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TransactionDetail), args.Error(1)
}

func (m *MockSubscriptionService) GetSubscriptionPlanByID(c *fiber.Ctx, planID uuid.UUID) (*model.SubscriptionPlan, error) {
	args := m.Called(c, planID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) UpdateSubscriptionPlan(c *fiber.Ctx, planID uuid.UUID, req *validation.UpdateSubscriptionPlan) (*model.SubscriptionPlan, error) {
	args := m.Called(c, planID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SubscriptionPlan), args.Error(1)
}

func generateTestToken(userID string) string {
	claims := jwt.MapClaims{
		"sub":  userID,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
		"type": config.TokenTypeAccess,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.JWTSecret))
	return tokenString
}

func TestFreemiumOrAccess(t *testing.T) {
	t.Run("should allow access with valid JWT and active freemium subscription", func(t *testing.T) {
		app := fiber.New()
		mockUserService := &MockUserService{}
		mockProductTokenService := &MockProductTokenService{}
		mockSubscriptionService := &MockSubscriptionService{}

		userID := uuid.New().String()
		user := fixture.UserWithFreemium()
		user.ID = uuid.MustParse(userID)

		subscription := &model.UserSubscriptionResponse{
			ID:       uuid.New(),
			UserID:   user.ID,
			Plan:     model.SubscriptionPlanResponse{Name: "Freemium Trial"},
			IsActive: true,
			EndDate:  time.Now().AddDate(0, 0, 14),
		}

		mockUserService.On("GetUserByID", mock.Anything, userID).Return(user, nil)
		mockSubscriptionService.On("GetUserActiveSubscription", mock.Anything, user.ID).Return(subscription, nil)

		token := generateTestToken(userID)

		app.Get("/test", middleware.FreemiumOrAccess(mockUserService, mockProductTokenService, mockSubscriptionService), func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		mockUserService.AssertExpectations(t)
		mockSubscriptionService.AssertExpectations(t)
		// ProductTokenService should not be called when user has active subscription
	})

	t.Run("should deny access without JWT token", func(t *testing.T) {
		app := fiber.New()
		mockUserService := &MockUserService{}
		mockProductTokenService := &MockProductTokenService{}
		mockSubscriptionService := &MockSubscriptionService{}

		app.Get("/test", middleware.FreemiumOrAccess(mockUserService, mockProductTokenService, mockSubscriptionService), func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("should deny access with expired freemium subscription", func(t *testing.T) {
		app := fiber.New()
		mockUserService := &MockUserService{}
		mockProductTokenService := &MockProductTokenService{}
		mockSubscriptionService := &MockSubscriptionService{}

		userID := uuid.New().String()
		user := fixture.UserWithExpiredFreemium()
		user.ID = uuid.MustParse(userID)

		subscription := &model.UserSubscriptionResponse{
			ID:       uuid.New(),
			UserID:   user.ID,
			Plan:     model.SubscriptionPlanResponse{Name: "Freemium Trial"},
			IsActive: true,
			EndDate:  time.Now().AddDate(0, 0, -1), // Expired yesterday
		}

		mockUserService.On("GetUserByID", mock.Anything, userID).Return(user, nil)
		mockSubscriptionService.On("GetUserActiveSubscription", mock.Anything, user.ID).Return(subscription, nil)
		mockProductTokenService.On("GetProductTokenByUserID", mock.Anything, user.ID).Return(nil, gorm.ErrRecordNotFound)

		token := generateTestToken(userID)

		app.Get("/test", middleware.FreemiumOrAccess(mockUserService, mockProductTokenService, mockSubscriptionService), func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "success"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode)

		mockUserService.AssertExpectations(t)
		mockSubscriptionService.AssertExpectations(t)
		mockProductTokenService.AssertExpectations(t)
	})
}
