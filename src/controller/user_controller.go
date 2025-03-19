package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController struct {
	UserService  service.UserService
	TokenService service.TokenService
}

func NewUserController(userService service.UserService, tokenService service.TokenService) *UserController {
	return &UserController{
		UserService:  userService,
		TokenService: tokenService,
	}
}

// @Tags         Users
// @Summary      Get all users
// @Description  Only admins can retrieve all users.
// @Security BearerAuth
// @Produce      json
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of users"    default(10)
// @Param        search   query     string  false  "Search by name or email or role"
// @Router       /users [get]
// @Success      200  {object}  example.GetAllUserResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
func (u *UserController) GetUsers(c *fiber.Ctx) error {
	query := &validation.QueryUser{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	users, totalResults, err := u.UserService.GetUsers(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
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

// @Tags         Users
// @Summary      Get a user
// @Description  Logged in users can fetch only their own user information. Only admins can fetch other users.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "User id"
// @Router       /users/{id} [get]
// @Success      200  {object}  example.GetUserResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (u *UserController) GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("userId")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := u.UserService.GetUserByID(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithUser{
			Status:  "success",
			Message: "Get user successfully",
			User:    *user,
		})
}

// @Tags         Users
// @Summary      Create a user
// @Description  Only admins can create other users.
// @Security BearerAuth
// @Produce      json
// @Param        request  body  validation.CreateUser  true  "Request body"
// @Router       /users [post]
// @Success      201  {object}  example.CreateUserResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      409  {object}  example.DuplicateEmail  "Email already taken"
func (u *UserController) CreateUser(c *fiber.Ctx) error {
	req := new(validation.CreateUser)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := u.UserService.CreateUser(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithUser{
			Status:  "success",
			Message: "Create user successfully",
			User:    *user,
		})
}

// @Tags         Users
// @Summary      Update a user
// @Description  Logged-in users can only update their own information. Only admins can update other users.
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        id               path      string    true   "User ID"
// @Param        name             formData  string    false  "User's name (max 50 characters)"
// @Param        email            formData  string    false  "User's email (must be valid email, max 50 characters)"
// @Param        password         formData  string    false  "Password (8-20 characters, must contain letters and numbers)"
// @Param        birth_date       formData  string    false  "Birth date (YYYY-MM-DD format)"
// @Param        height           formData  number    false  "Height in cm (0-300)"
// @Param        weight           formData  number    false  "Weight in kg (0-500)"
// @Param        gender           formData  string    false  "Gender (Male or Female)"
// @Param        activity_level   formData  string    false  "Activity level (Light, Medium, Heavy)"
// @Param        medical_history  formData  string    false  "Medical history (max 1000 characters)"
// @Param        profile_picture  formData  file      false  "Profile picture (optional)"
// @Router       /users/{id} [patch]
// @Success      200  {object}  example.UpdateUserResponse
func (u *UserController) UpdateUser(c *fiber.Ctx) error {
	req := new(validation.UpdateUser)
	userID := c.Params("userId")
	var profilePicture *multipart.FileHeader
	var hasFile bool

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid form data")
	}

	file, err := c.FormFile("profile_picture")
	if err == nil {
		profilePicture = file
		hasFile = true
	}

	user, err := u.UserService.UpdateUser(c, req, userID, profilePicture, hasFile)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithUser{
			Status:  "success",
			Message: "Update user successfully",
			User:    *user,
		})
}

// @Tags         Users
// @Summary      Delete a user
// @Description  Logged in users can delete only themselves. Only admins can delete other users.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "User id"
// @Router       /users/{id} [delete]
// @Success      200  {object}  example.DeleteUserResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (u *UserController) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("userId")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := u.TokenService.DeleteAllToken(c, userID); err != nil {
		return err
	}

	if err := u.UserService.DeleteUser(c, userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Status:  "success",
			Message: "Delete user successfully",
		})
}

// @Tags         Users
// @Summary      Get user statistics
// @Description  Get user's weight, height, and calorie statistics.
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "User ID"
// @Router       /users/{id}/statistics [get]
// @Success      200  {object}  example.UserStatisticsResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (u *UserController) GetUserStatistics(c *fiber.Ctx) error {
	userID := c.Params("userId")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	statistics, err := u.UserService.GetUserStatistics(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(statistics)
}
