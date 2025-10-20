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

	meal.Get("/", m.FreemiumOrAccess(u, p, ss), m.SubscriptionRequired(ss, "health_info"), mealController.GetMeals)
	meal.Post("/", m.FreemiumOrAccess(u, p, ss), mealController.AddMeal)
	meal.Post("/scan", m.FreemiumOrAccess(u, p, ss), mealController.ScanMeal)
	meal.Get("/:mealId", m.FreemiumOrAccess(u, p, ss), mealController.GetMealByID)
	meal.Put("/:mealId", m.FreemiumOrAccess(u, p, ss), mealController.UpdateMeal)
	meal.Delete("/:mealId", m.FreemiumOrAccess(u, p, ss), mealController.DeleteMeal)
	meal.Get("/:mealId/scan-detail", m.FreemiumOrAccess(u, p, ss), mealController.GetMealScanDetailByID)
	meal.Post("/:mealId/scan-detail", m.FreemiumOrAccess(u, p, ss), mealController.AddMealScanDetail)
}
