package service

import (
	"app/src/model"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ArticlesService interface {
	// Article methods
	CreateArticle(ctx *fiber.Ctx, article *model.Article) (*model.Article, error)
	GetArticles(ctx *fiber.Ctx) ([]model.Article, error)
	GetArticleByID(ctx *fiber.Ctx, articleID string) (*model.Article, error)
	UpdateArticle(ctx *fiber.Ctx, articleID string, article *model.Article) (*model.Article, error)
	DeleteArticle(ctx *fiber.Ctx, articleID string) error

	// Article Category methods
	CreateArticleCategory(ctx *fiber.Ctx, category *model.ArticleCategory) (*model.ArticleCategory, error)
	GetArticleCategories(ctx *fiber.Ctx) ([]model.ArticleCategory, error)
	DeleteArticleCategory(ctx *fiber.Ctx, categoryID string) error
}

type articlesService struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewArticlesService(db *gorm.DB) ArticlesService {
	return &articlesService{
		Log: logrus.New(),
		DB:  db,
	}
}

func (s *articlesService) CreateArticle(ctx *fiber.Ctx, article *model.Article) (*model.Article, error) {
	if err := s.DB.WithContext(ctx.Context()).Create(article).Error; err != nil {
		s.Log.Errorf("Failed to create article: %+v", err)
		return nil, err
	}
	return article, nil
}

func (s *articlesService) GetArticles(ctx *fiber.Ctx) ([]model.Article, error) {
	var articles []model.Article
	if err := s.DB.WithContext(ctx.Context()).
		Order("created_at DESC").
		Find(&articles).Error; err != nil {
		s.Log.Errorf("Failed to get articles: %+v", err)
		return nil, err
	}
	return articles, nil
}

func (s *articlesService) GetArticleByID(ctx *fiber.Ctx, articleID string) (*model.Article, error) {
	var article model.Article
	if err := s.DB.WithContext(ctx.Context()).
		Where("id = ?", articleID).
		First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Article not found")
		}
		s.Log.Errorf("Failed to get article: %+v", err)
		return nil, err
	}
	return &article, nil
}

func (s *articlesService) UpdateArticle(ctx *fiber.Ctx, articleID string, article *model.Article) (*model.Article, error) {
	existingArticle := new(model.Article)
	if err := s.DB.WithContext(ctx.Context()).
		Where("id = ?", articleID).
		First(existingArticle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Article not found")
		}
		s.Log.Errorf("Failed to find article: %+v", err)
		return nil, err
	}

	updates := make(map[string]interface{})
	if article.Title != "" {
		updates["title"] = article.Title
	}
	if article.CategoryID != nil {
		updates["category_id"] = article.CategoryID
	}
	if article.Slug != "" {
		updates["slug"] = article.Slug
	}
	if article.Image != nil {
		updates["image"] = article.Image
	}
	if article.Content != "" {
		updates["content"] = article.Content
	}
	if article.PublishedAt != nil {
		updates["published_at"] = article.PublishedAt
	}

	if len(updates) == 0 {
		return existingArticle, nil
	}

	if err := s.DB.WithContext(ctx.Context()).
		Model(existingArticle).
		Updates(updates).Error; err != nil {
		s.Log.Errorf("Failed to update article: %+v", err)
		return nil, err
	}

	return existingArticle, nil
}

func (s *articlesService) DeleteArticle(ctx *fiber.Ctx, articleID string) error {
	if err := s.DB.WithContext(ctx.Context()).
		Where("id = ?", articleID).
		Delete(&model.Article{}).Error; err != nil {
		s.Log.Errorf("Failed to delete article: %+v", err)
		return err
	}
	return nil
}

func (s *articlesService) CreateArticleCategory(ctx *fiber.Ctx, category *model.ArticleCategory) (*model.ArticleCategory, error) {
	if err := s.DB.WithContext(ctx.Context()).Create(category).Error; err != nil {
		s.Log.Errorf("Failed to create article category: %+v", err)
		return nil, err
	}
	return category, nil
}

func (s *articlesService) GetArticleCategories(ctx *fiber.Ctx) ([]model.ArticleCategory, error) {
	var categories []model.ArticleCategory
	if err := s.DB.WithContext(ctx.Context()).
		Order("created_at DESC").
		Find(&categories).Error; err != nil {
		s.Log.Errorf("Failed to get article categories: %+v", err)
		return nil, err
	}
	return categories, nil
}

func (s *articlesService) DeleteArticleCategory(ctx *fiber.Ctx, categoryID string) error {
	// First check if any articles are using this category
	var count int64
	if err := s.DB.WithContext(ctx.Context()).
		Model(&model.Article{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error; err != nil {
		s.Log.Errorf("Failed to check article category usage: %+v", err)
		return err
	}

	if count > 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot delete category that is in use by articles")
	}

	if err := s.DB.WithContext(ctx.Context()).
		Where("id = ?", categoryID).
		Delete(&model.ArticleCategory{}).Error; err != nil {
		s.Log.Errorf("Failed to delete article category: %+v", err)
		return err
	}
	return nil
}
