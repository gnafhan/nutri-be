package middleware

import (
	"app/src/config"
	"app/src/utils"
	"time"

	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
)

// SentryConfig returns the Sentry middleware configuration
func SentryConfig() fiber.Handler {
	// Initialize Sentry if not already initialized
	if sentry.CurrentHub().Client().Options().Dsn == "" {
		initSentry()
	}

	// Create Sentry handler with options
	sentryHandler := sentryfiber.New(sentryfiber.Options{
		Repanic:         true,
		WaitForDelivery: true,
	})

	return sentryHandler
}

// SentryEnhancement adds custom tags and context to Sentry events
func SentryEnhancement() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if hub := sentryfiber.GetHubFromContext(c); hub != nil {
			// Add request-specific tags
			hub.Scope().SetTag("method", c.Method())
			hub.Scope().SetTag("path", c.Path())
			hub.Scope().SetTag("user_agent", c.Get("User-Agent"))
			hub.Scope().SetTag("ip", c.IP())

			// Add environment context
			hub.Scope().SetTag("environment", config.SentryEnvironment)
			hub.Scope().SetTag("app_version", "1.0.0") // You can make this configurable

			// Add user context if available
			if userID := c.Locals("user_id"); userID != nil {
				hub.Scope().SetUser(sentry.User{
					ID: userID.(string),
				})
			}

			// Add request data as extra context
			hub.Scope().SetExtra("request_id", c.Locals("request_id"))
			hub.Scope().SetExtra("request_size", len(c.Body()))
		}

		return c.Next()
	}
}

// initSentry initializes Sentry with configuration
func initSentry() {
	// Use DSN from configuration
	dsn := config.SentryDSN
	if dsn == "" {
		utils.Log.Warn("SENTRY_DSN not configured, Sentry will not capture errors")
		return
	}

	// Set default environment if not provided
	environment := config.SentryEnvironment
	if environment == "" {
		if config.IsProd {
			environment = "production"
		} else {
			environment = "development"
		}
	}

	// Initialize Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: environment,
		Debug:       config.SentryDebug,
		// Set sample rate for performance monitoring
		TracesSampleRate: 0.1,
		// Set sample rate for profiling (removed as it's not available in this version)
		// Set release version
		Release: "nutribox-api@1.0.0",
		// Set server name
		ServerName: "nutribox-api",
		// Set before send hook to filter sensitive data
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Filter out sensitive data
			if event.Request != nil {
				// Remove sensitive headers
				if event.Request.Headers != nil {
					delete(event.Request.Headers, "Authorization")
					delete(event.Request.Headers, "Cookie")
					delete(event.Request.Headers, "X-API-Key")
				}
			}
			return event
		},
	})

	if err != nil {
		utils.Log.Errorf("Sentry initialization failed: %v", err)
	} else {
		utils.Log.Info("Sentry initialized successfully")
	}
}

// CaptureError captures an error with additional context
func CaptureError(err error, tags map[string]string, extra map[string]interface{}) {
	hub := sentry.CurrentHub()

	// Add tags
	for key, value := range tags {
		hub.Scope().SetTag(key, value)
	}

	// Add extra data
	for key, value := range extra {
		hub.Scope().SetExtra(key, value)
	}

	// Capture the error
	hub.CaptureException(err)
}

// CaptureMessage captures a message with additional context
func CaptureMessage(message string, level sentry.Level, tags map[string]string, extra map[string]interface{}) {
	hub := sentry.CurrentHub()

	// Add tags
	for key, value := range tags {
		hub.Scope().SetTag(key, value)
	}

	// Add extra data
	for key, value := range extra {
		hub.Scope().SetExtra(key, value)
	}

	// Capture the message
	hub.CaptureMessage(message)
}

// FlushSentry flushes any pending events before shutdown
func FlushSentry() {
	sentry.Flush(2 * time.Second)
}
