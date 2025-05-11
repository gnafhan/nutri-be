package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(v1 fiber.Router, userService service.UserService, tokenService service.TokenService, productTokenService service.ProductTokenService) {
	adminProductTokenController := controller.NewAdminProductTokenController(productTokenService)
	adminUserController := controller.NewAdminUserController(userService, tokenService)

	admin := v1.Group("/admin", m.Auth(userService, productTokenService))

	// Product Token routes
	productTokens := admin.Group("/product-tokens", m.Auth(userService, productTokenService, "getProductTokens"))
	productTokens.Get("/", adminProductTokenController.GetAllProductTokens)
	productTokens.Post("/", m.Auth(userService, productTokenService, "createProductToken"), adminProductTokenController.CreateProductToken)
	productTokens.Delete("/:id", m.Auth(userService, productTokenService, "deleteProductToken"), adminProductTokenController.DeleteProductToken)

	// User management routes
	users := admin.Group("/users", m.Auth(userService, productTokenService, "getUsers"))
	users.Get("/", adminUserController.GetAllUsers)
	users.Get("/:id", m.Auth(userService, productTokenService, "getUserDetails"), adminUserController.GetUserDetails)
	users.Patch("/:id", m.Auth(userService, productTokenService, "updateUser"), adminUserController.UpdateUser)
}
