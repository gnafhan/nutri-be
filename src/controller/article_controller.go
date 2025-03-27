package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ArticleController struct {
	ArticleService service.ArticlesService
}

func NewArticleController(service service.ArticlesService) *ArticleController {
	return &ArticleController{
		ArticleService: service,
	}
}

// @Tags         Articles
// @Summary      Create new article
// @Description  Create new article
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      example.CreateArticleRequest  true  "Article data"
// @Router       /articles [post]
// @Success      201  {object}  response.SuccessWithArticle
func (c *ArticleController) CreateArticle(ctx *fiber.Ctx) error {
	var request model.Article
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	user := ctx.Locals("user").(*model.User)
	request.UserID = user.ID

	article, err := c.ArticleService.CreateArticle(ctx, &request)
	if err != nil {
		return err
	}

	// After creating the article, fetch it with category info
	articleResponse, err := c.ArticleService.GetArticleByID(ctx, article.ID.String())
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.SuccessWithArticle{
		Status:  "success",
		Message: "Article created successfully",
		Data:    *articleResponse,
	})
}

// @Tags         Articles
// @Summary      Get all articles
// @Description  Get all articles
// @Security     BearerAuth
// @Produce      json
// @Router       /articles [get]
// @Success      200  {object}  response.SuccessWithArticleList
func (c *ArticleController) GetArticles(ctx *fiber.Ctx) error {
	articles, err := c.ArticleService.GetArticles(ctx)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithArticleList{
		Status:  "success",
		Message: "Articles fetched successfully",
		Data:    articles, // No conversion needed, already ArticleResponse
	})
}

// @Tags         Articles
// @Summary      Get article by ID
// @Description  Get article by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Article ID"
// @Router       /articles/{id} [get]
// @Success      200  {object}  response.SuccessWithArticle
func (c *ArticleController) GetArticleByID(ctx *fiber.Ctx) error {
	articleID := ctx.Params("id")
	if _, err := uuid.Parse(articleID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid article ID",
		})
	}

	article, err := c.ArticleService.GetArticleByID(ctx, articleID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithArticle{
		Status:  "success",
		Message: "Article fetched successfully",
		Data:    *article,
	})
}

// @Tags         Articles
// @Summary      Update article
// @Description  Update article
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Article ID"
// @Param        request  body      example.UpdateArticleRequest  true  "Article data"
// @Router       /articles/{id} [put]
// @Success      200  {object}  response.SuccessWithArticle
func (c *ArticleController) UpdateArticle(ctx *fiber.Ctx) error {
	articleID := ctx.Params("id")
	if _, err := uuid.Parse(articleID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid article ID",
		})
	}

	var request model.Article
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	user := ctx.Locals("user").(*model.User)
	request.UserID = user.ID

	if _, err := c.ArticleService.UpdateArticle(ctx, articleID, &request); err != nil {
		return err
	}

	// After updating, fetch the article with category info
	articleResponse, err := c.ArticleService.GetArticleByID(ctx, articleID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithArticle{
		Status:  "success",
		Message: "Article updated successfully",
		Data:    *articleResponse,
	})
}

// @Tags         Articles
// @Summary      Delete article
// @Description  Delete article
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Article ID"
// @Router       /articles/{id} [delete]
// @Success      200  {object}  response.Common
func (c *ArticleController) DeleteArticle(ctx *fiber.Ctx) error {
	articleID := ctx.Params("id")
	if _, err := uuid.Parse(articleID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid article ID",
		})
	}

	if err := c.ArticleService.DeleteArticle(ctx, articleID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Common{
		Status:  "success",
		Message: "Article deleted successfully",
	})
}

// Categories endpoints
// @Tags         Article Categories
// @Summary      Create new article category
// @Description  Create new article category
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      example.CreateArticleCategoryRequest  true  "Category data"
// @Router       /article-categories [post]
// @Success      201  {object}  response.SuccessWithArticleCategory
func (c *ArticleController) CreateArticleCategory(ctx *fiber.Ctx) error {
	var request model.ArticleCategory
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	user := ctx.Locals("user").(*model.User)
	request.UserID = user.ID

	category, err := c.ArticleService.CreateArticleCategory(ctx, &request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.SuccessWithArticleCategory{
		Status:  "success",
		Message: "Article category created successfully",
		Data:    *category,
	})
}

// @Tags         Article Categories
// @Summary      Get all article categories
// @Description  Get all article categories
// @Security     BearerAuth
// @Produce      json
// @Router       /article-categories [get]
// @Success      200  {object}  response.SuccessWithArticleCategoryList
func (c *ArticleController) GetArticleCategories(ctx *fiber.Ctx) error {
	categories, err := c.ArticleService.GetArticleCategories(ctx)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithArticleCategoryList{
		Status:  "success",
		Message: "Article categories fetched successfully",
		Data:    categories,
	})
}

// @Tags         Article Categories
// @Summary      Delete article category
// @Description  Delete article category
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Category ID"
// @Router       /article-categories/{id} [delete]
// @Success      200  {object}  response.Common
func (c *ArticleController) DeleteArticleCategory(ctx *fiber.Ctx) error {
	categoryID := ctx.Params("id")
	if _, err := uuid.Parse(categoryID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid category ID",
		})
	}
	if err := c.ArticleService.DeleteArticleCategory(ctx, categoryID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Common{
		Status:  "success",
		Message: "Article category deleted successfully",
	})
}
