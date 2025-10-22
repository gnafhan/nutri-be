package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func MealRoutes(v1 fiber.Router, u service.UserService, ml service.MealService, ss service.SubscriptionService) {
	mealController := controller.NewMealController(ml)

	meal := v1.Group("/meals")

	meal.Get("/", m.FreemiumOrAccess(u, nil, ss), m.SubscriptionRequired(ss, "health_info"), mealController.GetMeals)
	meal.Post("/", m.FreemiumOrAccess(u, nil, ss), mealController.AddMeal)
	meal.Post("/scan", m.FreemiumOrAccess(u, nil, ss), mealController.ScanMeal)
	meal.Get("/:mealId", m.FreemiumOrAccess(u, nil, ss), mealController.GetMealByID)
	meal.Put("/:mealId", m.FreemiumOrAccess(u, nil, ss), mealController.UpdateMeal)
	meal.Delete("/:mealId", m.FreemiumOrAccess(u, nil, ss), mealController.DeleteMeal)
	meal.Get("/:mealId/scan-detail", m.FreemiumOrAccess(u, nil, ss), mealController.GetMealScanDetailByID)
	meal.Post("/:mealId/scan-detail", m.FreemiumOrAccess(u, nil, ss), mealController.AddMealScanDetail)
}
