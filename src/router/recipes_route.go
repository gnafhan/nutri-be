package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func RecipeRoutes(v1 fiber.Router, u service.UserService, p service.ProductTokenService, recipeService service.RecipesService) {
	recipeController := controller.NewRecipesController(recipeService)

	recipes := v1.Group("/recipes")
	recipes.Get("/", m.Auth(u, p), recipeController.GetRecipes)
	recipes.Post("/", m.Auth(u, p, "manageUsers"), recipeController.CreateRecipe)
	recipes.Get("/:id", m.Auth(u, p), recipeController.GetRecipeByID)
	recipes.Put("/:id", m.Auth(u, p, "manageUsers"), recipeController.UpdateRecipe)
	recipes.Delete("/:id", m.Auth(u, p, "manageUsers"), recipeController.DeleteRecipe)
}
