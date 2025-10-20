package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func ArticleRoutes(v1 fiber.Router, u service.UserService, p service.ProductTokenService, ss service.SubscriptionService, articleService service.ArticlesService) {
	articleController := controller.NewArticleController(articleService)

	articles := v1.Group("/articles")
	articles.Get("/", m.FreemiumOrAccess(u, p, ss), articleController.GetArticles)
	articles.Post("/", m.Auth(u, p, "manageUsers"), articleController.CreateArticle)
	articles.Get("/:id", m.FreemiumOrAccess(u, p, ss), articleController.GetArticleByID)
	articles.Put("/:id", m.Auth(u, p, "manageUsers"), articleController.UpdateArticle)
	articles.Delete("/:id", m.Auth(u, p, "manageUsers"), articleController.DeleteArticle)

	categories := v1.Group("/article-categories")
	categories.Get("/", m.FreemiumOrAccess(u, p, ss), articleController.GetArticleCategories)
	categories.Post("/", m.Auth(u, p, "manageUsers"), articleController.CreateArticleCategory)
	categories.Delete("/:id", m.Auth(u, p, "manageUsers"), articleController.DeleteArticleCategory)
}
