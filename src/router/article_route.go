package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func ArticleRoutes(v1 fiber.Router, u service.UserService, p service.ProductTokenService, articleService service.ArticlesService) {
	articleController := controller.NewArticleController(articleService)

	articles := v1.Group("/articles")
	articles.Get("/", m.Auth(u, p), articleController.GetArticles)
	articles.Post("/", m.Auth(u, p), articleController.CreateArticle)
	articles.Get("/:id", m.Auth(u, p), articleController.GetArticleByID)
	articles.Put("/:id", m.Auth(u, p), articleController.UpdateArticle)
	articles.Delete("/:id", m.Auth(u, p), articleController.DeleteArticle)

	categories := v1.Group("/article-categories")
	categories.Get("/", m.Auth(u, p), articleController.GetArticleCategories)
	categories.Post("/", m.Auth(u, p), articleController.CreateArticleCategory)
	categories.Delete("/:id", m.Auth(u, p), articleController.DeleteArticleCategory)
}
