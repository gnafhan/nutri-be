package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func LoginStreakRoutes(
	router fiber.Router,
	userService service.UserService,
	productTokenService service.ProductTokenService,
	loginStreakService service.LoginStreakService,
) {
	loginStreakController := controller.NewLoginStreakController(loginStreakService)

	loginStreakRouter := router.Group("/login-streak")

	// Middleware to verify user token
	loginStreakRouter.Use(m.Auth(userService, productTokenService))

	// Record login streak (increments streak when user opens app)
	loginStreakRouter.Post("/record", loginStreakController.RecordLoginStreak)

	// Get login streak information
	loginStreakRouter.Get("/", loginStreakController.GetLoginStreak)
}
