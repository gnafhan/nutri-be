package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func SubscriptionRoutes(
	v1 fiber.Router,
	u service.UserService,
	p service.ProductTokenService,
	subService service.SubscriptionService,
) {
	subController := controller.NewSubscriptionController(subService)

	subGroup := v1.Group("/subscriptions")
	{
		subGroup.Get("/plans", subController.GetPlans)

		// Webhook endpoint for payment notification - doesn't require auth
		subGroup.Post("/notification", subController.HandlePaymentNotification)
		// Also handle the route with trailing slash
		subGroup.Post("/notification/", subController.HandlePaymentNotification)

		// Authenticated endpoints
		authGroup := subGroup.Group("", m.Auth(u, p))
		{
			authGroup.Get("/me", subController.GetMySubscription)
			authGroup.Get("/check-feature", subController.CheckFeatureAccess)
			authGroup.Post("/purchase/:planID", subController.PurchasePlan)

		}
	}
}
