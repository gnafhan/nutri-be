package service

import (
	"app/src/model"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RecipesService interface {
	CreateRecipe(ctx *fiber.Ctx, recipe *model.Recipe) (*model.Recipe, error)
	GetRecipes(ctx *fiber.Ctx) ([]model.Recipe, error)
	GetRecipeByID(ctx *fiber.Ctx, recipeID string) (*model.Recipe, error)
	UpdateRecipe(ctx *fiber.Ctx, recipeID string, recipe *model.Recipe) (*model.Recipe, error)
	DeleteRecipe(ctx *fiber.Ctx, recipeID string) error
}

type recipesService struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewRecipesService(db *gorm.DB) RecipesService {
	return &recipesService{
		Log: logrus.New(),
		DB:  db,
	}
}

func (s *recipesService) CreateRecipe(ctx *fiber.Ctx, recipe *model.Recipe) (*model.Recipe, error) {
	if err := s.DB.WithContext(ctx.Context()).Create(recipe).Error; err != nil {
		s.Log.Errorf("Failed to create recipe: %+v", err)
		return nil, err
	}
	return recipe, nil
}

func (s *recipesService) GetRecipes(ctx *fiber.Ctx) ([]model.Recipe, error) {
	var recipes []model.Recipe
	if err := s.DB.WithContext(ctx.Context()).
		Order("created_at DESC").
		Find(&recipes).Error; err != nil {
		s.Log.Errorf("Failed to get recipes: %+v", err)
		return nil, err
	}
	return recipes, nil
}

func (s *recipesService) GetRecipeByID(ctx *fiber.Ctx, recipeID string) (*model.Recipe, error) {
	var recipe model.Recipe
	if err := s.DB.WithContext(ctx.Context()).
		Where("id = ?", recipeID).
		First(&recipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Recipe not found")
		}
		s.Log.Errorf("Failed to get recipe: %+v", err)
		return nil, err
	}
	return &recipe, nil
}

func (s *recipesService) UpdateRecipe(ctx *fiber.Ctx, recipeID string, recipe *model.Recipe) (*model.Recipe, error) {
	existingRecipe := new(model.Recipe)
	if err := s.DB.WithContext(ctx.Context()).
		Where("id = ?", recipeID).
		First(existingRecipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Recipe not found")
		}
		s.Log.Errorf("Failed to find recipe: %+v", err)
		return nil, err
	}

	updates := make(map[string]interface{})
	if recipe.Name != "" {
		updates["name"] = recipe.Name
	}
	if recipe.Slug != "" {
		updates["slug"] = recipe.Slug
	}
	if recipe.Image != nil {
		updates["image"] = recipe.Image
	}
	if recipe.Description != "" {
		updates["description"] = recipe.Description
	}
	if recipe.Ingredients != "" {
		updates["ingredients"] = recipe.Ingredients
	}
	if recipe.Instructions != "" {
		updates["instructions"] = recipe.Instructions
	}
	if recipe.Label != nil {
		updates["label"] = recipe.Label
	}

	if len(updates) == 0 {
		return existingRecipe, nil
	}

	if err := s.DB.WithContext(ctx.Context()).
		Model(existingRecipe).
		Updates(updates).Error; err != nil {
		s.Log.Errorf("Failed to update recipe: %+v", err)
		return nil, err
	}

	return existingRecipe, nil
}

func (s *recipesService) DeleteRecipe(ctx *fiber.Ctx, recipeID string) error {
	if err := s.DB.WithContext(ctx.Context()).
		Where("id = ?", recipeID).
		Delete(&model.Recipe{}).Error; err != nil {
		s.Log.Errorf("Failed to delete recipe: %+v", err)
		return err
	}
	return nil
}
