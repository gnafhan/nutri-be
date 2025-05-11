package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func MealRoutes(v1 fiber.Router, u service.UserService, p service.ProductTokenService, ml service.MealService, ss service.SubscriptionService) {
	mealController := controller.NewMealController(ml)

	meal := v1.Group("/meals")

	meal.Get("/", m.Auth(u, p), m.SubscriptionRequired(ss, "health_info"), mealController.GetMeals)
	meal.Post("/", m.Auth(u, p), mealController.AddMeal)
	meal.Post("/scan", m.Auth(u, p), mealController.ScanMeal)
	meal.Get("/:mealId", m.Auth(u, p), mealController.GetMealByID)
	meal.Put("/:mealId", m.Auth(u, p), mealController.UpdateMeal)
	meal.Delete("/:mealId", m.Auth(u, p), mealController.DeleteMeal)
	meal.Get("/:mealId/scan-detail", m.Auth(u, p), mealController.GetMealScanDetailByID)
	meal.Post("/:mealId/scan-detail", m.Auth(u, p), mealController.AddMealScanDetail)
}
