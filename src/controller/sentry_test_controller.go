package controller

import (
	"app/src/utils"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SentryTestController handles Sentry testing endpoints
type SentryTestController struct{}

// NewSentryTestController creates a new Sentry test controller
func NewSentryTestController() *SentryTestController {
	return &SentryTestController{}
}

// TestError tests Sentry error reporting
// @Summary Test Sentry error reporting
// @Description This endpoint intentionally throws an error to test Sentry integration
// @Tags Sentry
// @Accept json
// @Produce json
// @Success 500 {object} response.ErrorResponse
// @Router /v1/sentry/test-error [get]
func (c *SentryTestController) TestError(ctx *fiber.Ctx) error {
	sentryService := utils.NewSentryService()

	// Add breadcrumb before error
	sentryService.AddBreadcrumb(ctx, "About to throw test error", "test", "info", map[string]interface{}{
		"test_type": "intentional_error",
		"timestamp": time.Now().UTC(),
	})

	// Capture a test error
	err := errors.New("This is a test error for Sentry integration")
	sentryService.CaptureErrorWithContext(ctx, err, map[string]string{
		"test_type": "sentry_integration_test",
		"endpoint":  "/v1/sentry/test-error",
	}, map[string]interface{}{
		"test_data": "This error was intentionally generated to test Sentry",
		"timestamp": time.Now().UTC(),
	})

	return ctx.Status(500).JSON(fiber.Map{
		"error":   "Test error generated for Sentry testing",
		"message": "Check your Sentry dashboard for the error report",
	})
}

// TestMessage tests Sentry message reporting
// @Summary Test Sentry message reporting
// @Description This endpoint sends a test message to Sentry
// @Tags Sentry
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /v1/sentry/test-message [get]
func (c *SentryTestController) TestMessage(ctx *fiber.Ctx) error {
	sentryService := utils.NewSentryService()

	// Add breadcrumb
	sentryService.AddBreadcrumb(ctx, "About to send test message", "test", "info", map[string]interface{}{
		"test_type": "intentional_message",
		"timestamp": time.Now().UTC(),
	})

	// Capture a test message
	sentryService.CaptureMessageWithContext(ctx, "This is a test message for Sentry integration", "info", map[string]string{
		"test_type": "sentry_integration_test",
		"endpoint":  "/v1/sentry/test-message",
	}, map[string]interface{}{
		"test_data": "This message was intentionally sent to test Sentry",
		"timestamp": time.Now().UTC(),
	})

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Test message sent to Sentry successfully",
	})
}

// TestPanic tests Sentry panic recovery
// @Summary Test Sentry panic recovery
// @Description This endpoint intentionally panics to test Sentry panic recovery
// @Tags Sentry
// @Accept json
// @Produce json
// @Success 500 {object} response.ErrorResponse
// @Router /v1/sentry/test-panic [get]
func (c *SentryTestController) TestPanic(ctx *fiber.Ctx) error {
	sentryService := utils.NewSentryService()

	// Add breadcrumb before panic
	sentryService.AddBreadcrumb(ctx, "About to panic for testing", "test", "error", map[string]interface{}{
		"test_type": "intentional_panic",
		"timestamp": time.Now().UTC(),
	})

	// This will trigger a panic, which should be caught by Sentry middleware
	panic("This is a test panic for Sentry integration")
}
