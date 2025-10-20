package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func UsersWeightHeightRoutes(v1 fiber.Router, u service.UserService, p service.ProductTokenService, ss service.SubscriptionService, uwhService service.UsersWeightHeightService) {
	uwhController := controller.NewUsersWeightHeightController(uwhService)

	uwh := v1.Group("/weight-height")

	uwh.Get("/target/", m.FreemiumOrAccess(u, p, ss), uwhController.GetWeightHeightsTarget)
	uwh.Post("/target/", m.FreemiumOrAccess(u, p, ss), uwhController.AddWeightHeightTarget)
	uwh.Get("/target/:uwhId", m.FreemiumOrAccess(u, p, ss), uwhController.GetWeightHeightTargetByID)
	uwh.Put("/target/:uwhId", m.FreemiumOrAccess(u, p, ss), uwhController.UpdateWeightHeightTarget)
	uwh.Delete("/target/:uwhId", m.FreemiumOrAccess(u, p, ss), uwhController.DeleteWeightHeightTarget)

	uwh.Get("/", m.FreemiumOrAccess(u, p, ss), uwhController.GetWeightHeights)
	uwh.Post("/", m.FreemiumOrAccess(u, p, ss), uwhController.AddWeightHeight)
	uwh.Get("/:uwhId", m.FreemiumOrAccess(u, p, ss), uwhController.GetWeightHeightByID)
	uwh.Put("/:uwhId", m.FreemiumOrAccess(u, p, ss), uwhController.UpdateWeightHeight)
	uwh.Delete("/:uwhId", m.FreemiumOrAccess(u, p, ss), uwhController.DeleteWeightHeight)
}
