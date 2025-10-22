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
	paymentService := service.NewMidtransPaymentService()
	subscriptionService := service.NewSubscriptionService(db, paymentService)
	userService := service.NewUserService(db, validate, subscriptionService)
	tokenService := service.NewTokenService(db, validate, userService, subscriptionService)
	authService := service.NewAuthService(db, validate, userService, tokenService, subscriptionService)
	mealService := service.NewMealService(db, config.LogMealApiKey, config.LogMealBaseUrl)
	uwhService := service.NewUsersWeightHeightService(db)
	articleService := service.NewArticlesService(db)
	recipesService := service.NewRecipesService(db)
	loginStreakService := service.NewLoginStreakService(db, validate)
	bahanMakananService := service.NewBahanMakananService(client)
	productTokenService := service.NewProductTokenService(db, validate)

	v1 := app.Group("/v1")

	HealthCheckRoutes(v1, healthCheckService)
	AuthRoutes(v1, authService, userService, tokenService, emailService)
	UserRoutes(v1, userService, tokenService)
	MealRoutes(v1, userService, mealService, subscriptionService)
	UsersWeightHeightRoutes(v1, userService, subscriptionService, uwhService)
	ArticleRoutes(v1, userService, subscriptionService, articleService)
	RecipeRoutes(v1, userService, subscriptionService, recipesService)
	SubscriptionRoutes(v1, userService, subscriptionService)
	ProductTokenRoutes(v1, userService, productTokenService)
	AdminRoutes(v1, userService, tokenService, subscriptionService)
	LoginStreakRoutes(v1, userService, subscriptionService, loginStreakService)
	BahanMakananRoutes(v1, userService, subscriptionService, bahanMakananService)
	HomeRoutes(v1, userService, subscriptionService, mealService)

	// TODO: add another routes here...

	if !config.IsProd {
		DocsRoutes(v1)
		SentryTestRoutes(v1) // Only add Sentry test routes in development
	}
}
