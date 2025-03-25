package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RecipesController struct {
	RecipeService service.RecipesService
}

func NewRecipesController(service service.RecipesService) *RecipesController {
	return &RecipesController{
		RecipeService: service,
	}
}

// @Tags         Recipes
// @Summary      Create new recipe
// @Description  Create new recipe
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      example.CreateRecipeRequest  true  "Recipe data"
// @Router       /recipes [post]
// @Success      201  {object}  response.SuccessWithRecipe
func (c *RecipesController) CreateRecipe(ctx *fiber.Ctx) error {
	var request model.Recipe
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	recipe, err := c.RecipeService.CreateRecipe(ctx, &request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.SuccessWithRecipe{
		Status:  "success",
		Message: "Recipe created successfully",
		Data:    *recipe,
	})
}

// @Tags         Recipes
// @Summary      Get all recipes
// @Description  Get all recipes
// @Security     BearerAuth
// @Produce      json
// @Router       /recipes [get]
// @Success      200  {object}  response.SuccessWithRecipeList
func (c *RecipesController) GetRecipes(ctx *fiber.Ctx) error {
	recipes, err := c.RecipeService.GetRecipes(ctx)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithRecipeList{
		Status:  "success",
		Message: "Recipes fetched successfully",
		Data:    recipes,
	})
}

// @Tags         Recipes
// @Summary      Get recipe by ID
// @Description  Get recipe by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Recipe ID"
// @Router       /recipes/{id} [get]
// @Success      200  {object}  response.SuccessWithRecipe
func (c *RecipesController) GetRecipeByID(ctx *fiber.Ctx) error {
	recipeID := ctx.Params("id")
	if _, err := uuid.Parse(recipeID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid recipe ID",
		})
	}

	recipe, err := c.RecipeService.GetRecipeByID(ctx, recipeID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithRecipe{
		Status:  "success",
		Message: "Recipe fetched successfully",
		Data:    *recipe,
	})
}

// @Tags         Recipes
// @Summary      Update recipe
// @Description  Update recipe
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Recipe ID"
// @Param        request  body      example.UpdateRecipeRequest  true  "Recipe data"
// @Router       /recipes/{id} [put]
// @Success      200  {object}  response.SuccessWithRecipe
func (c *RecipesController) UpdateRecipe(ctx *fiber.Ctx) error {
	recipeID := ctx.Params("id")
	if _, err := uuid.Parse(recipeID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid recipe ID",
		})
	}

	var request model.Recipe
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	recipe, err := c.RecipeService.UpdateRecipe(ctx, recipeID, &request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithRecipe{
		Status:  "success",
		Message: "Recipe updated successfully",
		Data:    *recipe,
	})
}

// @Tags         Recipes
// @Summary      Delete recipe
// @Description  Delete recipe
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Recipe ID"
// @Router       /recipes/{id} [delete]
// @Success      200  {object}  response.Common
func (c *RecipesController) DeleteRecipe(ctx *fiber.Ctx) error {
	recipeID := ctx.Params("id")
	if _, err := uuid.Parse(recipeID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Status:  "error",
			Message: "Invalid recipe ID",
		})
	}

	if err := c.RecipeService.DeleteRecipe(ctx, recipeID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Common{
		Status:  "success",
		Message: "Recipe deleted successfully",
	})
}
