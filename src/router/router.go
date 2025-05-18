package router

import (
	"app/src/config"
	"app/src/grpc"
	"app/src/service"
	"app/src/validation"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()
	grpcServerAddr := fmt.Sprintf("%s:%s", config.GRPC_HOST, config.GRPC_PORT)
	client, _ := grpc.NewBahanMakananClient(grpcServerAddr)

	healthCheckService := service.NewHealthCheckService(db)
	emailService := service.NewEmailService()
	userService := service.NewUserService(db, validate)
	paymentService := service.NewMidtransPaymentService()
	subscriptionService := service.NewSubscriptionService(db, paymentService)
	tokenService := service.NewTokenService(db, validate, userService, subscriptionService)
	authService := service.NewAuthService(db, validate, userService, tokenService)
	productTokenService := service.NewProductTokenService(db, validate)
	mealService := service.NewMealService(db, config.LogMealApiKey, config.LogMealBaseUrl)
	uwhService := service.NewUsersWeightHeightService(db)
	articleService := service.NewArticlesService(db)
	recipesService := service.NewRecipesService(db)
	loginStreakService := service.NewLoginStreakService(db, validate)
	bahanMakananService := service.NewBahanMakananService(client)

	v1 := app.Group("/v1")

	HealthCheckRoutes(v1, healthCheckService)
	AuthRoutes(v1, authService, userService, productTokenService, tokenService, emailService)
	UserRoutes(v1, userService, productTokenService, tokenService)
	ProductTokenRoutes(v1, userService, productTokenService)
	MealRoutes(v1, userService, productTokenService, mealService, subscriptionService)
	UsersWeightHeightRoutes(v1, userService, productTokenService, uwhService)
	ArticleRoutes(v1, userService, productTokenService, articleService)
	RecipeRoutes(v1, userService, productTokenService, recipesService)
	SubscriptionRoutes(v1, userService, productTokenService, subscriptionService)
	AdminRoutes(v1, userService, tokenService, productTokenService, subscriptionService)
	LoginStreakRoutes(v1, userService, productTokenService, loginStreakService)
	BahanMakananRoutes(v1, userService, productTokenService, bahanMakananService)

	// TODO: add another routes here...

	if !config.IsProd {
		DocsRoutes(v1)
	}
}
