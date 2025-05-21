package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminProductTokenController struct {
	ProductTokenService service.ProductTokenService
}

func NewAdminProductTokenController(
	productTokenService service.ProductTokenService,
) *AdminProductTokenController {
	return &AdminProductTokenController{
		ProductTokenService: productTokenService,
	}
}

// @Tags         Admin
// @Summary      Get all product tokens
// @Description  Returns a list of all product tokens with their activation status
// @Produce      json
// @Security     BearerAuth
// @Param        with_user   query     boolean  false  "Include user data"
// @Router       /admin/product-tokens [get]
// @Success      200  {object}  example.GetAllProductTokensResponse
// @Failure      403  {object}  response.ErrorResponse
func (c *AdminProductTokenController) GetAllProductTokens(ctx *fiber.Ctx) error {
	query := &validation.ProductTokenQuery{
		WithUser: ctx.QueryBool("with_user", false),
	}

	tokens, err := c.ProductTokenService.GetAllProductTokens(ctx, query)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithProductTokens{
		Status:  "success",
		Message: "Product tokens retrieved successfully",
		Data:    tokens,
	})
}

// @Tags         Admin
// @Summary      Create new product token
// @Description  Creates a new custom product token that can be used by users. Optionally, a `subscription_plan_id` can be provided to grant a subscription when the token is verified.
// @Accept       json
// @Produce      json
// @Param        request   body  validation.CreateCustomToken  true  "Product token details (token, is_active, subscription_plan_id)"
// @Security     BearerAuth
// @Router       /admin/product-tokens [post]
// @Success      201  {object}  example.CreateProductTokenResponse
// @Failure      400  {object}  response.ErrorResponse "Invalid request or token already exists"
// @Failure      403  {object}  response.ErrorResponse
func (c *AdminProductTokenController) CreateProductToken(ctx *fiber.Ctx) error {
	req := new(validation.CreateCustomToken)

	// Default value for IsActive
	req.IsActive = true

	// Jika ada body, parse
	if ctx.Get("Content-Type") == "application/json" {
		if err := ctx.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}
	}

	token, err := c.ProductTokenService.CreateProductToken(ctx, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.SuccessWithProductToken{
		Status:  "success",
		Message: "Product token created successfully",
		Data:    *token,
	})
}

// @Tags         Admin
// @Summary      Delete product token
// @Description  Deletes a product token by ID
// @Produce      json
// @Param        id   path  string  true  "Product token ID"
// @Security     BearerAuth
// @Router       /admin/product-tokens/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminProductTokenController) DeleteProductToken(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	if idParam == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Product token ID is required")
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid product token ID format")
	}

	if err := c.ProductTokenService.AdminDeleteProductToken(ctx, id); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Common{
		Status:  "success",
		Message: "Product token deleted successfully",
	})
}

// @Tags         Admin
// @Summary      Update product token
// @Description  Updates an existing product token. Fields to update (token, is_active, subscription_plan_id) should be provided in the request body. To remove a subscription plan, pass an empty string for `subscription_plan_id`.
// @Accept       json
// @Produce      json
// @Param        id   path  string  true  "Product token ID"
// @Param        request   body  validation.UpdateProductToken  true  "Product token details to update"
// @Security     BearerAuth
// @Router       /admin/product-tokens/{id} [put]
// @Success      200  {object}  response.SuccessWithProductToken
// @Failure      400  {object}  response.ErrorResponse "Invalid request, token already exists, or invalid SubscriptionPlanID"
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse "Product token or Subscription Plan not found"
func (c *AdminProductTokenController) UpdateProductToken(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	if idParam == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Product token ID is required")
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid product token ID format")
	}

	req := new(validation.UpdateProductToken)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	// Call the service to update the product token
	updatedToken, err := c.ProductTokenService.UpdateProductToken(ctx, id, req)
	if err != nil {
		return err // Return the error from the service directly (it should be a fiber.Error)
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithProductToken{
		Status:  "success",
		Message: "Product token updated successfully",
		Data:    *updatedToken,
	})
}
