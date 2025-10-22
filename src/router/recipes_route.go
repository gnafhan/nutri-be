package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func RecipeRoutes(v1 fiber.Router, u service.UserService, ss service.SubscriptionService, recipeService service.RecipesService) {
	recipeController := controller.NewRecipesController(recipeService)

	recipes := v1.Group("/recipes")
	recipes.Get("/", m.FreemiumOrAccess(u, nil, ss), recipeController.GetRecipes)
	recipes.Post("/", m.Auth(u, nil, "manageUsers"), recipeController.CreateRecipe)
	recipes.Get("/:id", m.FreemiumOrAccess(u, nil, ss), recipeController.GetRecipeByID)
	recipes.Put("/:id", m.Auth(u, nil, "manageUsers"), recipeController.UpdateRecipe)
	recipes.Delete("/:id", m.Auth(u, nil, "manageUsers"), recipeController.DeleteRecipe)
}
