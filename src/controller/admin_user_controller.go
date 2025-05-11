package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminUserController struct {
	UserService  service.UserService
	TokenService service.TokenService
}

func NewAdminUserController(
	userService service.UserService,
	tokenService service.TokenService,
) *AdminUserController {
	return &AdminUserController{
		UserService:  userService,
		TokenService: tokenService,
	}
}

// @Tags         Admin
// @Summary      Get all users
// @Description  Admin endpoint to retrieve all users with pagination
// @Security     BearerAuth
// @Produce      json
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of users"    default(10)
// @Param        search   query     string  false  "Search by name or email or role"
// @Router       /admin/users [get]
// @Success      200  {object}  example.SuccessWithPaginateUsers
// @Failure      403  {object}  response.ErrorResponse
func (c *AdminUserController) GetAllUsers(ctx *fiber.Ctx) error {
	query := &validation.QueryUser{
		Page:   ctx.QueryInt("page", 1),
		Limit:  ctx.QueryInt("limit", 10),
		Search: ctx.Query("search", ""),
	}

	users, totalResults, err := c.UserService.GetUsers(ctx, query)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).
		JSON(response.SuccessWithPaginate[model.User]{
			Status:       "success",
			Message:      "Get all users successfully",
			Results:      users,
			Page:         query.Page,
			Limit:        query.Limit,
			TotalPages:   int64(math.Ceil(float64(totalResults) / float64(query.Limit))),
			TotalResults: totalResults,
		})
}

// @Tags         Admin
// @Summary      Get user details
// @Description  Admin endpoint to get detailed user information
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "User id"
// @Router       /admin/users/{id} [get]
// @Success      200  {object}  response.SuccessWithUser
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminUserController) GetUserDetails(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := c.UserService.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).
		JSON(response.SuccessWithUser{
			Status:  "success",
			Message: "Get user details successfully",
			User:    *user,
		})
}

// @Tags         Admin
// @Summary      Update user
// @Description  Admin endpoint to update user information
// @Security     BearerAuth
// @Accept       application/json
// @Produce      json
// @Param        id       path      string    true   "User ID"
// @Param        request  body      validation.UpdateUser  true  "Update user data"
// @Router       /admin/users/{id} [patch]
// @Success      200  {object}  response.SuccessWithUser
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminUserController) UpdateUser(ctx *fiber.Ctx) error {
	req := new(validation.UpdateUser)
	userID := ctx.Params("id")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := c.UserService.UpdateUser(ctx, req, userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).
		JSON(response.SuccessWithUser{
			Status:  "success",
			Message: "Update user successfully",
			User:    *user,
		})
}
