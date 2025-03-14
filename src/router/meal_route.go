package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func MealRoutes(v1 fiber.Router, u service.UserService, p service.ProductTokenService, ml service.MealService) {
	mealController := controller.NewMealController(ml)

	meal := v1.Group("/meals")

	meal.Post("/scan", m.Auth(u, p), mealController.ScanMeal)
	// meal.Post("/history", m.Auth(u, p), mealController.GetMealHistory)
	// meal.Get("/history/:mealId", m.Auth(u, p), mealController.GetMealHistoryDetail)
}
