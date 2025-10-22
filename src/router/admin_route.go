package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(v1 fiber.Router, userService service.UserService, tokenService service.TokenService, subscriptionService service.SubscriptionService) {
	adminUserController := controller.NewAdminUserController(userService, tokenService)
	adminSubscriptionController := controller.NewAdminSubscriptionController(subscriptionService)

	admin := v1.Group("/admin", m.Auth(userService, nil))

	// User management routes
	users := admin.Group("/users", m.Auth(userService, nil, "getUsers"))
	users.Get("/", adminUserController.GetAllUsers)
	users.Get("/:id", m.Auth(userService, nil, "getUserDetails"), adminUserController.GetUserDetails)
	users.Patch("/:id", m.Auth(userService, nil, "updateUser"), adminUserController.UpdateUser)

	// Subscription routes
	subscriptions := admin.Group("/subscriptions", m.Auth(userService, nil, "getSubscriptions"))
	subscriptions.Get("/", adminSubscriptionController.GetAllUserSubscriptions)

	// Specific subscription routes
	subscription := subscriptions.Group("/:subscription_id")
	subscription.Get("/", adminSubscriptionController.GetUserSubscriptionDetails)
	subscription.Patch("/", adminSubscriptionController.UpdateUserSubscription, m.Auth(userService, nil, "manageSubscriptions"))
	subscription.Get("/transactions", adminSubscriptionController.GetTransactionLogs, m.Auth(userService, nil, "viewTransactions"))
	subscription.Patch("/payment-status", adminSubscriptionController.UpdatePaymentStatus, m.Auth(userService, nil, "updatePaymentStatus"))

	// Subscription plans routes
	subscriptionPlans := admin.Group("/subscription-plans", m.Auth(userService, nil, "getSubscriptionPlans"))
	subscriptionPlans.Get("/", adminSubscriptionController.GetAllSubscriptionPlans)
	subscriptionPlans.Get("/:plan_id", adminSubscriptionController.GetSubscriptionPlanByID)
	subscriptionPlans.Patch("/:plan_id", adminSubscriptionController.UpdateSubscriptionPlan, m.Auth(userService, nil, "manageSubscriptionPlans"))

	// All transactions route
	transactions := admin.Group("/transactions", m.Auth(userService, nil, "viewTransactions"))
	transactions.Get("/", adminSubscriptionController.GetAllTransactions)
	transactions.Get("/:id", adminSubscriptionController.GetTransactionByID)
}
