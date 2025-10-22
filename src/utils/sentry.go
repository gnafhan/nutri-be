package utils

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
)

// SentryService provides utility functions for Sentry integration
type SentryService struct{}

// NewSentryService creates a new Sentry service instance
func NewSentryService() *SentryService {
	return &SentryService{}
}

// CaptureErrorWithContext captures an error with Fiber context
func (s *SentryService) CaptureErrorWithContext(c *fiber.Ctx, err error, tags map[string]string, extra map[string]interface{}) {
	hub := sentryfiber.GetHubFromContext(c)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	// Add request context
	if c != nil {
		hub.Scope().SetTag("method", c.Method())
		hub.Scope().SetTag("path", c.Path())
		hub.Scope().SetTag("user_agent", c.Get("User-Agent"))
		hub.Scope().SetTag("ip", c.IP())
		hub.Scope().SetTag("status_code", fmt.Sprintf("%d", c.Response().StatusCode()))

		// Add request data
		hub.Scope().SetExtra("request_body", string(c.Body()))
		hub.Scope().SetExtra("query_params", c.Queries())
		hub.Scope().SetExtra("headers", c.GetReqHeaders())
	}

	// Add custom tags
	for key, value := range tags {
		hub.Scope().SetTag(key, value)
	}

	// Add custom extra data
	for key, value := range extra {
		hub.Scope().SetExtra(key, value)
	}

	// Add user context if available
	if userID := c.Locals("user_id"); userID != nil {
		hub.Scope().SetUser(sentry.User{
			ID: userID.(string),
		})
	}

	// Capture the error
	hub.CaptureException(err)
}

// CaptureMessageWithContext captures a message with Fiber context
func (s *SentryService) CaptureMessageWithContext(c *fiber.Ctx, message string, level sentry.Level, tags map[string]string, extra map[string]interface{}) {
	hub := sentryfiber.GetHubFromContext(c)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	// Add request context
	if c != nil {
		hub.Scope().SetTag("method", c.Method())
		hub.Scope().SetTag("path", c.Path())
		hub.Scope().SetTag("user_agent", c.Get("User-Agent"))
		hub.Scope().SetTag("ip", c.IP())
		hub.Scope().SetTag("status_code", fmt.Sprintf("%d", c.Response().StatusCode()))
	}

	// Add custom tags
	for key, value := range tags {
		hub.Scope().SetTag(key, value)
	}

	// Add custom extra data
	for key, value := range extra {
		hub.Scope().SetExtra(key, value)
	}

	// Add user context if available
	if userID := c.Locals("user_id"); userID != nil {
		hub.Scope().SetUser(sentry.User{
			ID: userID.(string),
		})
	}

	// Set the level and capture the message
	hub.Scope().SetLevel(level)
	hub.CaptureMessage(message)
}

// CaptureAPIError captures API-specific errors with enhanced context
func (s *SentryService) CaptureAPIError(c *fiber.Ctx, err error, endpoint string, operation string) {
	tags := map[string]string{
		"error_type":  "api_error",
		"endpoint":    endpoint,
		"operation":   operation,
		"environment": "development", // Will be set by middleware
	}

	extra := map[string]interface{}{
		"timestamp":   time.Now().UTC(),
		"app_version": "1.0.0",
	}

	s.CaptureErrorWithContext(c, err, tags, extra)
}

// CaptureDatabaseError captures database-specific errors
func (s *SentryService) CaptureDatabaseError(c *fiber.Ctx, err error, operation string, table string) {
	tags := map[string]string{
		"error_type":  "database_error",
		"operation":   operation,
		"table":       table,
		"environment": "development", // Will be set by middleware
	}

	extra := map[string]interface{}{
		"timestamp":   time.Now().UTC(),
		"app_version": "1.0.0",
	}

	s.CaptureErrorWithContext(c, err, tags, extra)
}

// CaptureValidationError captures validation errors
func (s *SentryService) CaptureValidationError(c *fiber.Ctx, err error, field string, value interface{}) {
	tags := map[string]string{
		"error_type":  "validation_error",
		"field":       field,
		"environment": "development", // Will be set by middleware
	}

	extra := map[string]interface{}{
		"timestamp":   time.Now().UTC(),
		"field_value": value,
		"app_version": "1.0.0",
	}

	s.CaptureErrorWithContext(c, err, tags, extra)
}

// CaptureExternalAPIError captures external API errors
func (s *SentryService) CaptureExternalAPIError(c *fiber.Ctx, err error, apiName string, endpoint string, statusCode int) {
	tags := map[string]string{
		"error_type":  "external_api_error",
		"api_name":    apiName,
		"endpoint":    endpoint,
		"status_code": string(rune(statusCode)),
		"environment": "development", // Will be set by middleware
	}

	extra := map[string]interface{}{
		"timestamp":         time.Now().UTC(),
		"api_name":          apiName,
		"external_endpoint": endpoint,
		"status_code":       statusCode,
		"app_version":       "1.0.0",
	}

	s.CaptureErrorWithContext(c, err, tags, extra)
}

// CapturePerformanceIssue captures performance-related issues
func (s *SentryService) CapturePerformanceIssue(c *fiber.Ctx, message string, duration time.Duration, operation string) {
	tags := map[string]string{
		"issue_type":  "performance",
		"operation":   operation,
		"environment": "development", // Will be set by middleware
	}

	extra := map[string]interface{}{
		"timestamp":   time.Now().UTC(),
		"duration_ms": duration.Milliseconds(),
		"operation":   operation,
		"app_version": "1.0.0",
	}

	s.CaptureMessageWithContext(c, message, sentry.LevelWarning, tags, extra)
}

// CaptureSecurityEvent captures security-related events
func (s *SentryService) CaptureSecurityEvent(c *fiber.Ctx, event string, details map[string]interface{}) {
	tags := map[string]string{
		"event_type":  "security",
		"environment": "development", // Will be set by middleware
	}

	extra := map[string]interface{}{
		"timestamp":   time.Now().UTC(),
		"app_version": "1.0.0",
	}

	// Add security event details
	for key, value := range details {
		extra[key] = value
	}

	s.CaptureMessageWithContext(c, event, sentry.LevelWarning, tags, extra)
}

// AddBreadcrumb adds a breadcrumb to the current scope
func (s *SentryService) AddBreadcrumb(c *fiber.Ctx, message string, category string, level sentry.Level, data map[string]interface{}) {
	hub := sentryfiber.GetHubFromContext(c)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	breadcrumb := &sentry.Breadcrumb{
		Message:   message,
		Category:  category,
		Level:     level,
		Data:      data,
		Timestamp: time.Now(),
	}

	hub.AddBreadcrumb(breadcrumb, &sentry.BreadcrumbHint{})
}

// SetUserContext sets user context for error tracking
func (s *SentryService) SetUserContext(c *fiber.Ctx, userID string, email string, username string) {
	hub := sentryfiber.GetHubFromContext(c)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.Scope().SetUser(sentry.User{
		ID:       userID,
		Email:    email,
		Username: username,
	})
}

// Flush flushes any pending Sentry events
func (s *SentryService) Flush() {
	sentry.Flush(2 * time.Second)
}
