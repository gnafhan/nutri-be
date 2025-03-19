package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UsersWeightHeightController struct {
	UsersWeightHeightService service.UsersWeightHeightService
}

func NewUsersWeightHeightController(service service.UsersWeightHeightService) *UsersWeightHeightController {
	return &UsersWeightHeightController{
		UsersWeightHeightService: service,
	}
}

// @Tags         Weight Height Record
// @Summary      Add a new weight and height record
// @Description  Logged in users can add a new weight and height record.
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      example.AddWeightHeightRequest  true  "Weight and height data"
// @Router       /weight-height [post]
// @Success      201  {object}  example.AddWeightHeightResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (c *UsersWeightHeightController) AddWeightHeight(ctx *fiber.Ctx) error {
	var request model.UsersWeightHeightHistory
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Common{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	user := ctx.Locals("user").(*model.User)
	request.UserID = user.ID

	result, err := c.UsersWeightHeightService.AddWeightHeight(ctx, &request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Common{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.SuccessWithWeightHeight{
		Status:  "success",
		Message: "Weight and height record added successfully",
		Data:    *result,
	})
}

// @Tags         Weight Height Record
// @Summary      Get all weight and height records
// @Description  Logged in users can fetch their own weight and height records.
// @Security     BearerAuth
// @Produce      json
// @Router       /weight-height [get]
// @Success      200  {object}  example.GetAllWeightHeightResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (c *UsersWeightHeightController) GetWeightHeights(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*model.User)

	records, err := c.UsersWeightHeightService.GetWeightHeights(ctx, user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Common{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithWeightHeightList{
		Status:  "success",
		Message: "Weight and height records fetched successfully",
		Data:    records,
	})
}

// @Tags         Weight Height Record
// @Summary      Get a weight height
// @Description  Logged in users can fetch their own weight and height record.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Record id"
// @Router       /weight-height/{id} [get]
// @Success      200  {object}  example.GetWeightHeightResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (c *UsersWeightHeightController) GetWeightHeightByID(ctx *fiber.Ctx) error {
	uwhId := ctx.Params("uwhId")
	user := ctx.Locals("user").(*model.User)

	if _, err := uuid.Parse(uwhId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid weight height ID")
	}

	uwh, err := c.UsersWeightHeightService.GetWeightHeightByID(ctx, uwhId, user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).
		JSON(response.SuccessWithWeightHeight{
			Status:  "success",
			Message: "Get weight height successfully",
			Data:    *uwh,
		})
}

// @Tags         Weight Height Record
// @Summary      Update a weight and height record
// @Description  Logged in users can update their own weight and height records.
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                        true  "Record ID"
// @Param        request  body      example.UpdateWeightHeightRequest  true  "Weight and height data"
// @Router       /weight-height/{id} [put]
// @Success      200  {object}  example.UpdateWeightHeightResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (c *UsersWeightHeightController) UpdateWeightHeight(ctx *fiber.Ctx) error {
	recordID := ctx.Params("uwhId")

	if _, err := uuid.Parse(recordID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Common{
			Status:  "error",
			Message: "Invalid record ID",
		})
	}

	var request model.UsersWeightHeightHistory
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Common{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	user := ctx.Locals("user").(*model.User)
	request.UserID = user.ID

	result, err := c.UsersWeightHeightService.UpdateWeightHeight(ctx, recordID, &request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithWeightHeight{
		Status:  "success",
		Message: "Weight and height record updated successfully",
		Data:    *result,
	})
}

// @Tags         Weight Height Record
// @Summary      Delete a weight and height record
// @Description  Logged in users can delete their own weight and height records.
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "Record ID"
// @Router       /weight-height/{id} [delete]
// @Success      200  {object}  example.DeleteWeightHeightResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (c *UsersWeightHeightController) DeleteWeightHeight(ctx *fiber.Ctx) error {
	recordID := ctx.Params("uwhId")

	if _, err := uuid.Parse(recordID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Common{
			Status:  "error",
			Message: "Invalid record ID",
		})
	}

	user := ctx.Locals("user").(*model.User)

	if err := c.UsersWeightHeightService.DeleteWeightHeight(ctx, recordID, user.ID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Common{
		Status:  "success",
		Message: "Weight and height record deleted successfully",
	})
}
