package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MealController struct {
	MealService service.MealService
}

func NewMealController(ms service.MealService) *MealController {
	return &MealController{
		MealService: ms,
	}
}

// @Tags         Meals
// @Summary      Scan a meal
// @Description  Only users who already logged in and had product token verified can scan a meal an get the nutritions
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file      true  "Meal's image"
// @Router       /meals/scan [post]
// @Success      200  {object}  example.MealScanResponse
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

// @Tags         Meals
// @Summary      Get a user's meals
// @Description  Logged in users can fetch only their own meals information.
// @Security BearerAuth
// @Produce      json
// @Param        page  query     int     false  "Page number"  default(1)
// @Param        limit query     int     false  "Maximum number of meals per page"  default(10)
// @Router       /meals [get]
// @Success      200  {object}  example.GetAllMealsResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (mc *MealController) GetMeals(c *fiber.Ctx) error {
	meals, totalResults, err := mc.MealService.GetMeals(c)
	if err != nil {
		return err
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	totalPages := int64(math.Ceil(float64(totalResults) / float64(limit)))

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithPaginate[model.MealHistory]{
			Status:       "success",
			Message:      "Get user's meals successfully",
			Results:      meals,
			Page:         page,
			Limit:        limit,
			TotalPages:   totalPages,
			TotalResults: totalResults,
		})
}

// @Tags         Meals
// @Summary      Get a meal
// @Description  Logged in users can fetch only their own meal detail information. Only admins can fetch other user's meal.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Meal id"
// @Router       /meals/{id} [get]
// @Success      200  {object}  example.GetMealResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (mc *MealController) GetMealByID(c *fiber.Ctx) error {
	mealId := c.Params("mealId")

	if _, err := uuid.Parse(mealId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid meal ID")
	}

	meal, err := mc.MealService.GetMealByID(c, mealId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithMeal{
			Status:  "success",
			Message: "Get meal successfully",
			Meal:    *meal,
		})
}

// @Tags         Meals
// @Summary      Get a meal's scan detail
// @Description  Logged in users can fetch only their own meal's scan detail detail information. Only admins can fetch other user's meal's scan detail.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Meal id"
// @Router       /meals/{id}/scan-detail [get]
// @Success      200  {object}  example.GetMealScanDetailResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (mc *MealController) GetMealScanDetailByID(c *fiber.Ctx) error {
	mealId := c.Params("mealId")

	if _, err := uuid.Parse(mealId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid meal ID")
	}

	mealScanDetail, err := mc.MealService.GetMealScanDetailByID(c, mealId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithMealScanDetail{
			Status:         "success",
			Message:        "Get meal's scan detail successfully",
			MealScanDetail: *mealScanDetail,
		})
}

// @Tags         Meals
// @Summary      Add a new meal's scan detail
// @Description  Logged in users can add a new meal's scan detail.
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        mealId   path      string                      true  "Meal ID"
// @Param        request  body      example.AddMealScanDetailRequest  true  "Meal's scan detail data"
// @Router       /meals/{mealId}/scan-detail [post]
// @Success      201  {object}  example.AddMealScanDetailResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      404  {object}  example.NotFound  "Not found"
func (mc *MealController) AddMealScanDetail(c *fiber.Ctx) error {
	mealId := c.Params("mealId")

	if _, err := uuid.Parse(mealId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid meal ID")
	}

	var request model.MealHistoryDetail
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Common{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	mealScanDetail, err := mc.MealService.AddMealScanDetail(c, mealId, &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Common{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithMealScanDetail{
		Status:         "success",
		Message:        "Meal's scan detail added successfully",
		MealScanDetail: *mealScanDetail,
	})
}

// @Tags         Meals
// @Summary      Add a new meal
// @Description  Logged in users can add a new meal.
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      example.AddMealRequest  true  "Meal data"
// @Router       /meals [post]
// @Success      201  {object}  example.AddMealResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (mc *MealController) AddMeal(c *fiber.Ctx) error {
	var request model.MealHistory
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Common{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	user := c.Locals("user")
	userData := user.(*model.User)
	request.UserID = userData.ID

	meal, err := mc.MealService.AddMeal(c, &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Common{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithMeal{
		Status:  "success",
		Message: "Meal added successfully",
		Meal:    *meal,
	})
}

// @Tags         Meals
// @Summary      Update a meal
// @Description  Logged in users can update their own meal. Only admins can update other user's meal.
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                  true  "Meal ID"
// @Param        request  body      example.UpdateMealRequest  true  "Meal data"
// @Router       /meals/{id} [put]
// @Success      200  {object}  example.UpdateMealResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (mc *MealController) UpdateMeal(c *fiber.Ctx) error {
	mealId := c.Params("mealId")

	if _, err := uuid.Parse(mealId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid meal ID")
	}

	var request model.MealHistory
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Common{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	meal, err := mc.MealService.UpdateMeal(c, mealId, &request)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithMeal{
		Status:  "success",
		Message: "Meal updated successfully",
		Meal:    *meal,
	})
}

// @Tags         Meals
// @Summary      Delete a meal
// @Description  Logged in users can delete their own meal. Only admins can delete other user's meal.
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Meal ID"
// @Router       /meals/{id} [delete]
// @Success      200  {object}  example.DeleteMealResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (mc *MealController) DeleteMeal(c *fiber.Ctx) error {
	mealId := c.Params("mealId")

	if _, err := uuid.Parse(mealId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid meal ID")
	}

	if err := mc.MealService.DeleteMeal(c, mealId); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Status:  "success",
		Message: "Meal deleted successfully",
	})
}

// @Tags         Statistics
// @Summary      Get home statistics
// @Description  Logged in users can fetch their home statistics including today's consumed calories and weight/height info
// @Security     BearerAuth
// @Produce      json
// @Router       /home/statistic [get]
// @Success      200  {object}  example.GetHomeStatisticsResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (mc *MealController) GetHomeStatistics(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.User)

	homeStats, err := mc.MealService.GetHomeStatistics(c, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Common{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithHomeStatistics{
		Status:  "success",
		Message: "Home statistics fetched successfully",
		Data:    *homeStats,
	})
}
