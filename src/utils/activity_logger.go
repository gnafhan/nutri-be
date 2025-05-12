package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	// ActivityLog is the logger for user activities
	ActivityLog *logrus.Logger
	// RequestLog is the logger for API requests and responses
	RequestLog *logrus.Logger
)

// ActivityData represents user activity data to be logged
type ActivityData struct {
	UserID      string      `json:"userID"`
	Action      string      `json:"action"`
	Resource    string      `json:"resource,omitempty"`
	ResourceID  string      `json:"resourceID,omitempty"`
	Details     interface{} `json:"details,omitempty"`
	RequestID   string      `json:"requestID"`
	IPAddress   string      `json:"ipAddress,omitempty"`
	UserAgent   string      `json:"userAgent,omitempty"`
	StatusCode  int         `json:"statusCode,omitempty"`
	ElapsedTime string      `json:"elapsedTime,omitempty"`
}

// RequestResponseData represents API request and response data to be logged
type RequestResponseData struct {
	Method      string      `json:"method"`
	Path        string      `json:"path"`
	RequestID   string      `json:"requestID"`
	IPAddress   string      `json:"ipAddress"`
	UserAgent   string      `json:"userAgent,omitempty"`
	RequestBody interface{} `json:"requestBody,omitempty"`
	StatusCode  int         `json:"statusCode"`
	Response    interface{} `json:"response,omitempty"`
	ElapsedTime string      `json:"elapsedTime"`
	UserID      string      `json:"userID,omitempty"`
}

func init() {
	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		Log.Fatalf("Failed to create log directory: %v", err)
	}

	// Initialize activity logger
	ActivityLog = logrus.New()
	activityLogFile, err := os.OpenFile(
		filepath.Join(logDir, "activity.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		Log.Fatalf("Failed to open activity log file: %v", err)
	}
	ActivityLog.SetOutput(activityLogFile)
	ActivityLog.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	ActivityLog.SetLevel(logrus.InfoLevel)

	// Initialize request logger
	RequestLog = logrus.New()
	requestLogFile, err := os.OpenFile(
		filepath.Join(logDir, "request.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		Log.Fatalf("Failed to open request log file: %v", err)
	}
	RequestLog.SetOutput(requestLogFile)
	RequestLog.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	RequestLog.SetLevel(logrus.InfoLevel)
}

// LogUserActivity logs user activity
func LogUserActivity(data ActivityData) {
	if data.RequestID == "" {
		data.RequestID = fmt.Sprintf("req-%s", uuid.New().String()[:8])
	}

	ActivityLog.WithFields(logrus.Fields{
		"userID":      data.UserID,
		"action":      data.Action,
		"resource":    data.Resource,
		"resourceID":  data.ResourceID,
		"details":     data.Details,
		"requestID":   data.RequestID,
		"ipAddress":   data.IPAddress,
		"userAgent":   data.UserAgent,
		"statusCode":  data.StatusCode,
		"elapsedTime": data.ElapsedTime,
	}).Info("User activity")
}

// LogAPIRequest logs API request
func LogAPIRequest(data RequestResponseData) {
	if data.RequestID == "" {
		data.RequestID = fmt.Sprintf("req-%s", uuid.New().String()[:8])
	}

	RequestLog.WithFields(logrus.Fields{
		"method":      data.Method,
		"path":        data.Path,
		"requestID":   data.RequestID,
		"ipAddress":   data.IPAddress,
		"userAgent":   data.UserAgent,
		"requestBody": data.RequestBody,
		"userID":      data.UserID,
	}).Info("API Request")
}

// LogAPIResponse logs API response
func LogAPIResponse(data RequestResponseData) {
	RequestLog.WithFields(logrus.Fields{
		"method":      data.Method,
		"path":        data.Path,
		"requestID":   data.RequestID,
		"statusCode":  data.StatusCode,
		"response":    data.Response,
		"elapsedTime": data.ElapsedTime,
		"userID":      data.UserID,
	}).Info("API Response")
}

// Helper functions for common user activities

// LogLogin logs user login activity
func LogLogin(c *fiber.Ctx, userID string, success bool) {
	requestID := getRequestID(c)
	LogUserActivity(ActivityData{
		UserID:     userID,
		Action:     "login",
		Details:    map[string]interface{}{"success": success},
		RequestID:  requestID,
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		StatusCode: c.Response().StatusCode(),
	})
}

// LogRegistration logs user registration activity
func LogRegistration(c *fiber.Ctx, userID string) {
	requestID := getRequestID(c)
	LogUserActivity(ActivityData{
		UserID:     userID,
		Action:     "register",
		RequestID:  requestID,
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		StatusCode: c.Response().StatusCode(),
	})
}

// LogSubscriptionPurchase logs subscription purchase activity
func LogSubscriptionPurchase(c *fiber.Ctx, userID string, planID string, paymentMethod string) {
	requestID := getRequestID(c)
	LogUserActivity(ActivityData{
		UserID:     userID,
		Action:     "subscription_purchase",
		Resource:   "subscription_plan",
		ResourceID: planID,
		Details: map[string]interface{}{
			"payment_method": paymentMethod,
		},
		RequestID:  requestID,
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		StatusCode: c.Response().StatusCode(),
	})
}

// LogScanActivity logs food scan activity
func LogScanActivity(c *fiber.Ctx, userID string, foodItem string, calories int) {
	requestID := getRequestID(c)
	LogUserActivity(ActivityData{
		UserID:   userID,
		Action:   "scan_food",
		Resource: "food_item",
		Details: map[string]interface{}{
			"food_item": foodItem,
			"calories":  calories,
		},
		RequestID:  requestID,
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		StatusCode: c.Response().StatusCode(),
	})
}

// LogMealTracking logs meal tracking activity
func LogMealTracking(c *fiber.Ctx, userID string, mealType string, mealID string) {
	requestID := getRequestID(c)
	LogUserActivity(ActivityData{
		UserID:     userID,
		Action:     "track_meal",
		Resource:   "meal",
		ResourceID: mealID,
		Details: map[string]interface{}{
			"meal_type": mealType,
		},
		RequestID:  requestID,
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		StatusCode: c.Response().StatusCode(),
	})
}

// LogWeightUpdate logs weight update activity
func LogWeightUpdate(c *fiber.Ctx, userID string, weight float64) {
	requestID := getRequestID(c)
	LogUserActivity(ActivityData{
		UserID:   userID,
		Action:   "update_weight",
		Resource: "weight_record",
		Details: map[string]interface{}{
			"weight": weight,
		},
		RequestID:  requestID,
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		StatusCode: c.Response().StatusCode(),
	})
}

// Helper function to get request ID from context
func getRequestID(c *fiber.Ctx) string {
	requestID := c.Locals("requestID")
	if requestID == nil {
		return fmt.Sprintf("req-%s", uuid.New().String()[:8])
	}
	return requestID.(string)
}
