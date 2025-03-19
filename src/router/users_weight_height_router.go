package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func UsersWeightHeightRoutes(v1 fiber.Router, u service.UserService, p service.ProductTokenService, uwhService service.UsersWeightHeightService) {
	uwhController := controller.NewUsersWeightHeightController(uwhService)

	uwh := v1.Group("/weight-height")

	uwh.Get("/", m.Auth(u, p), uwhController.GetWeightHeights)
	uwh.Post("/", m.Auth(u, p), uwhController.AddWeightHeight)
	uwh.Get("/:uwhId", m.Auth(u, p), uwhController.GetWeightHeightByID)
	uwh.Put("/:uwhId", m.Auth(u, p), uwhController.UpdateWeightHeight)
	uwh.Delete("/:uwhId", m.Auth(u, p), uwhController.DeleteWeightHeight)
}
