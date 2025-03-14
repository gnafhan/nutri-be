package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func ProductTokenRoutes(v1 fiber.Router, u service.UserService, p service.ProductTokenService) {
	productTokenController := controller.NewProductTokenController(p)

	productToken := v1.Group("/product-token")

	productToken.Post("/verify", m.AuthWithoutTokenCheck(u), productTokenController.VerifyProductToken)
}
