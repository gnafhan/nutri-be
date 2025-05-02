package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/response/example"
	"app/src/service"
	"app/src/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SubscriptionController struct {
	Service service.SubscriptionService
}

func NewSubscriptionController(service service.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{
		Service: service,
	}
}

// @Tags         Subscription
// @Summary      Get all subscription plans
// @Description  Get available subscription plans
// @Produce      json
// @Router       /subscriptions/plans [get]
// @Success      200  {object}  response.SubscriptionPlansResponse
// @Success      200  {object}  example.SubscriptionPlanResponse
func (c *SubscriptionController) GetPlans(ctx *fiber.Ctx) error {
	plans, err := c.Service.GetAllPlans(ctx)
	if err != nil {
		return utils.APIError(ctx, fiber.StatusInternalServerError, "server_error", err.Error())
	}

	return ctx.JSON(response.SubscriptionPlansResponse{
		Status:  "success",
		Message: "Subscription plans retrieved",
		Data:    plans,
	})
}

// @Tags         Subscription
// @Summary      Purchase subscription plan
// @Description  Purchase a subscription plan
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        planID  path  string  true  "Plan ID"
// @Param        request  body  example.PurchaseSubscriptionRequest  true  "Payment data"
// @Router       /subscriptions/purchase/{planID} [post]
// @Success      200  {object}  response.UserSubscriptionResponse
// @Success      200  {object}  example.UserSubscriptionResponse
func (c *SubscriptionController) PurchasePlan(ctx *fiber.Ctx) error {
	planID := ctx.Params("planID")
	if _, err := uuid.Parse(planID); err != nil {
		return utils.APIError(ctx, fiber.StatusBadRequest, "invalid_id", "Invalid plan ID")
	}

	var req example.PurchaseSubscriptionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.APIError(ctx, fiber.StatusBadRequest, "invalid_request", err.Error())
	}

	user := ctx.Locals("user").(*model.User)
	subscription, err := c.Service.PurchasePlan(ctx, user.ID, uuid.MustParse(planID), req.PaymentMethod)
	if err != nil {
		return utils.APIError(ctx, fiber.StatusInternalServerError, "purchase_failed", err.Error())
	}

	return ctx.JSON(response.UserSubscriptionResponse{
		Status:  "success",
		Message: "Subscription purchased successfully",
		Data:    *subscription,
	})
}

// @Tags         Subscription
// @Summary      Get current subscription
// @Description  Get user's active subscription
// @Security     BearerAuth
// @Produce      json
// @Router       /subscriptions/me [get]
// @Success      200  {object}  response.UserSubscriptionResponse
// @Success      200  {object}  example.UserSubscriptionResponse
func (c *SubscriptionController) GetMySubscription(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*model.User)
	subscription, err := c.Service.GetUserActiveSubscription(ctx, user.ID)
	if err != nil {
		return utils.APIError(ctx, fiber.StatusNotFound, "no_subscription", "No active subscription found")
	}

	return ctx.JSON(response.UserSubscriptionResponse{
		Status:  "success",
		Message: "Active subscription retrieved",
		Data:    *subscription,
	})
}

// @Tags         Subscription
// @Summary      Check feature access
// @Description  Check if user has access to a feature
// @Security     BearerAuth
// @Produce      json
// @Param        feature  query  string  true  "Feature name"
// @Router       /subscriptions/check-feature [get]
// @Success      200  {object}  response.FeatureAccessResponse
func (c *SubscriptionController) CheckFeatureAccess(ctx *fiber.Ctx) error {
	feature := ctx.Query("feature")
	if feature == "" {
		return utils.APIError(ctx, fiber.StatusBadRequest, "missing_feature", "Feature parameter is required")
	}

	user := ctx.Locals("user").(*model.User)
	hasAccess, err := c.Service.CheckFeatureAccess(ctx, user.ID, feature)
	if err != nil {
		return utils.APIError(ctx, fiber.StatusInternalServerError, "check_failed", err.Error())
	}

	return ctx.JSON(response.FeatureAccessResponse{
		Status:  "success",
		Message: "Feature access checked",
		Data: response.FeatureData{
			Feature: feature,
			Access:  hasAccess,
		},
	})
}
