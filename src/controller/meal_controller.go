package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

type MealController struct {
	MealService service.MealService
}

func NewMealController(ms service.MealService) *MealController {
	return &MealController{
		MealService: ms,
	}
}

func (mc *MealController) ScanMeal(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Common{
			Status:  "error",
			Message: "Image file is required",
		})
	}

	user := c.Locals("user")
	userData := user.(*model.User)

	result, err := mc.MealService.ScanMeal(c, file, userData.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Common{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
