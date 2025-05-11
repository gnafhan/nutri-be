package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(v1 fiber.Router, userService service.UserService, tokenService service.TokenService, productTokenService service.ProductTokenService, subscriptionService service.SubscriptionService) {
	adminProductTokenController := controller.NewAdminProductTokenController(productTokenService)
	adminUserController := controller.NewAdminUserController(userService, tokenService)
	adminSubscriptionController := controller.NewAdminSubscriptionController(subscriptionService)

	admin := v1.Group("/admin", m.Auth(userService, productTokenService))

	// Product Token routes
	productTokens := admin.Group("/product-tokens", m.Auth(userService, productTokenService, "getProductTokens"))
	productTokens.Get("/", adminProductTokenController.GetAllProductTokens)
	productTokens.Post("/", m.Auth(userService, productTokenService, "createProductToken"), adminProductTokenController.CreateProductToken)
	productTokens.Delete("/:id", m.Auth(userService, productTokenService, "deleteProductToken"), adminProductTokenController.DeleteProductToken)

	// User management routes
	users := admin.Group("/users", m.Auth(userService, productTokenService, "getUsers"))
	users.Get("/", adminUserController.GetAllUsers)
	users.Get("/:id", m.Auth(userService, productTokenService, "getUserDetails"), adminUserController.GetUserDetails)
	users.Patch("/:id", m.Auth(userService, productTokenService, "updateUser"), adminUserController.UpdateUser)

	// Subscription routes
	subscriptions := admin.Group("/subscriptions", m.Auth(userService, productTokenService, "getSubscriptions"))
	subscriptions.Get("/", adminSubscriptionController.GetAllUserSubscriptions)

	// Specific subscription routes
	subscription := subscriptions.Group("/:subscription_id")
	subscription.Get("/", adminSubscriptionController.GetUserSubscriptionDetails)
	subscription.Patch("/", adminSubscriptionController.UpdateUserSubscription, m.Auth(userService, productTokenService, "manageSubscriptions"))
	subscription.Get("/transactions", adminSubscriptionController.GetTransactionLogs, m.Auth(userService, productTokenService, "viewTransactions"))
	subscription.Patch("/payment-status", adminSubscriptionController.UpdatePaymentStatus, m.Auth(userService, productTokenService, "updatePaymentStatus"))

	// Subscription plans routes
	subscriptionPlans := admin.Group("/subscription-plans", m.Auth(userService, productTokenService, "getSubscriptionPlans"))
	subscriptionPlans.Get("/", adminSubscriptionController.GetAllSubscriptionPlans)

	// All transactions route
	transactions := admin.Group("/transactions", m.Auth(userService, productTokenService, "viewTransactions"))
	transactions.Get("/", adminSubscriptionController.GetAllTransactions)
	transactions.Get("/:id", adminSubscriptionController.GetTransactionByID)
}
