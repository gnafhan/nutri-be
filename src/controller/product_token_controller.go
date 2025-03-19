package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

type ProductTokenController struct {
	ProductTokenService service.ProductTokenService
}

func NewProductTokenController(
	productTokenService service.ProductTokenService,
) *ProductTokenController {
	return &ProductTokenController{
		ProductTokenService: productTokenService,
	}
}

// @Tags         Product Token
// @Summary      Verify Product Token
// @Produce      json
// @Param        token   query  string  true  "The product token"
// @Security     BearerAuth
// @Router       /product-token/verify [post]
// @Success      200  {object}  example.VerifyProductTokenResponse
// @Failure      404  {object}  example.FailedVerifyProductToken  "Invalid or already used product token"
func (p *ProductTokenController) VerifyProductToken(c *fiber.Ctx) error {
	query := &validation.Token{
		Token: c.Query("token"),
	}

	if err := p.ProductTokenService.VerifyProductToken(c, query); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Status:  "success",
			Message: "Verify product token successfully",
		})
}
