package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func UsersWeightHeightRoutes(v1 fiber.Router, u service.UserService, ss service.SubscriptionService, uwhService service.UsersWeightHeightService) {
	uwhController := controller.NewUsersWeightHeightController(uwhService)

	uwh := v1.Group("/weight-height")

	uwh.Get("/target/", m.FreemiumOrAccess(u, nil, ss), uwhController.GetWeightHeightsTarget)
	uwh.Post("/target/", m.FreemiumOrAccess(u, nil, ss), uwhController.AddWeightHeightTarget)
	uwh.Get("/target/:uwhId", m.FreemiumOrAccess(u, nil, ss), uwhController.GetWeightHeightTargetByID)
	uwh.Put("/target/:uwhId", m.FreemiumOrAccess(u, nil, ss), uwhController.UpdateWeightHeightTarget)
	uwh.Delete("/target/:uwhId", m.FreemiumOrAccess(u, nil, ss), uwhController.DeleteWeightHeightTarget)

	uwh.Get("/", m.FreemiumOrAccess(u, nil, ss), uwhController.GetWeightHeights)
	uwh.Post("/", m.FreemiumOrAccess(u, nil, ss), uwhController.AddWeightHeight)
	uwh.Get("/:uwhId", m.FreemiumOrAccess(u, nil, ss), uwhController.GetWeightHeightByID)
	uwh.Put("/:uwhId", m.FreemiumOrAccess(u, nil, ss), uwhController.UpdateWeightHeight)
	uwh.Delete("/:uwhId", m.FreemiumOrAccess(u, nil, ss), uwhController.DeleteWeightHeight)
}
