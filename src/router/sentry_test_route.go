package router

import (
	"app/src/controller"

	"github.com/gofiber/fiber/v2"
)

// SentryTestRoutes sets up Sentry testing routes
func SentryTestRoutes(api fiber.Router) {
	sentryTestController := controller.NewSentryTestController()

	// Sentry testing routes (only in development)
	sentry := api.Group("/sentry")
	sentry.Get("/test-error", sentryTestController.TestError)
	sentry.Get("/test-message", sentryTestController.TestMessage)
	sentry.Get("/test-panic", sentryTestController.TestPanic)
}
