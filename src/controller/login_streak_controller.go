package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

// LoginStreakController handles login streak operations
type LoginStreakController struct {
	loginStreakService service.LoginStreakService
}

// LoginStreakResponse represents the response for login streak data
// @Description Response for login streak data
type LoginStreakResponse struct {
	Status        string                 `json:"status" example:"success"`
	Message       string                 `json:"message" example:"Login streak retrieved successfully"`
	CurrentStreak int                    `json:"current_streak" example:"5"`
	LongestStreak int                    `json:"longest_streak" example:"12"`
	WeeklyStreak  []LoginStreakDayDetail `json:"weekly_streak"`
}

// LoginStreakDayDetail represents streak info for a specific day of the week
// @Description Streak info for a specific day of the week
type LoginStreakDayDetail struct {
	DayOfWeek int    `json:"day_of_week" example:"1"` // 1 = Monday, ..., 7 = Sunday
	Date      string `json:"date" example:"2023-06-05T00:00:00Z"`
	HasLogin  bool   `json:"has_login" example:"true"`
}

func NewLoginStreakController(loginStreakService service.LoginStreakService) *LoginStreakController {
	return &LoginStreakController{
		loginStreakService: loginStreakService,
	}
}

// RecordLoginStreak records a login streak for the authenticated user
// @Summary Record login streak
// @Description Records a login streak for the authenticated user
// @Tags Login Streak
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.CommonResponse "Login streak recorded successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /v1/login-streak/record [post]
func (c *LoginStreakController) RecordLoginStreak(ctx *fiber.Ctx) error {
	// Get the user from context
	user := ctx.Locals("user").(*model.User)
	userID := user.ID

	err := c.loginStreakService.RecordLogin(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.CommonResponse{
		Status:  "success",
		Message: "Login streak recorded successfully",
	})
}

// GetLoginStreak retrieves the login streak information for the authenticated user
// @Summary Get login streak
// @Description Retrieves the login streak information for the authenticated user
// @Tags Login Streak
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessWithLoginStreak "Login streak data"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /v1/login-streak [get]
func (c *LoginStreakController) GetLoginStreak(ctx *fiber.Ctx) error {
	// Get the user from context
	user := ctx.Locals("user").(*model.User)
	userID := user.ID

	streakData, err := c.loginStreakService.GetLoginStreak(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Convert the model response to the expected response format
	responseData := model.LoginStreakData(*streakData)

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithLoginStreak{
		Status:  "success",
		Message: "Login streak retrieved successfully",
		Data:    responseData,
	})
}
