package router

import (
	"app/src/config"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()

	healthCheckService := service.NewHealthCheckService(db)
	emailService := service.NewEmailService()
	userService := service.NewUserService(db, validate)
	tokenService := service.NewTokenService(db, validate, userService)
	authService := service.NewAuthService(db, validate, userService, tokenService)
	productTokenService := service.NewProductTokenService(db, validate)
	mealService := service.NewMealService(db, config.LogMealApiKey, config.LogMealBaseUrl)
	uwhService := service.NewUsersWeightHeightService(db)
	articleService := service.NewArticlesService(db)
	recipesService := service.NewRecipesService(db)

	v1 := app.Group("/v1")

	HealthCheckRoutes(v1, healthCheckService)
	AuthRoutes(v1, authService, userService, productTokenService, tokenService, emailService)
	UserRoutes(v1, userService, productTokenService, tokenService)
	ProductTokenRoutes(v1, userService, productTokenService)
	MealRoutes(v1, userService, productTokenService, mealService)
	UsersWeightHeightRoutes(v1, userService, productTokenService, uwhService)
	ArticleRoutes(v1, userService, productTokenService, articleService)
	RecipeRoutes(v1, userService, productTokenService, recipesService)
	// TODO: add another routes here...

	if !config.IsProd {
		DocsRoutes(v1)
	}
}
