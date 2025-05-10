package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/utils"

	"encoding/json"

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
// @Param        request  body  model.PurchaseSubscriptionRequest  false  "Payment data (optional)"
// @Router       /subscriptions/purchase/{planID} [post]
// @Success      200  {object}  response.PaymentResponse
func (c *SubscriptionController) PurchasePlan(ctx *fiber.Ctx) error {
	planID := ctx.Params("planID")
	if _, err := uuid.Parse(planID); err != nil {
		return utils.APIError(ctx, fiber.StatusBadRequest, "invalid_id", "Invalid plan ID")
	}

	var req model.PurchaseSubscriptionRequest
	if err := ctx.BodyParser(&req); err != nil && err.Error() != "Unprocessable Entity" {
		return utils.APIError(ctx, fiber.StatusBadRequest, "invalid_request", err.Error())
	}

	user := ctx.Locals("user").(*model.User)
	paymentResponse, err := c.Service.PurchasePlan(ctx, user.ID, uuid.MustParse(planID), req.PaymentMethod)
	if err != nil {
		return utils.APIError(ctx, fiber.StatusInternalServerError, "purchase_failed", err.Error())
	}

	return ctx.JSON(response.PaymentResponse{
		Status:  "success",
		Message: "Payment initiated successfully",
		Data:    paymentResponse,
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

// @Tags         Subscription
// @Summary      Midtrans payment notification webhook
// @Description  Handle payment notification from Midtrans
// @Accept       json
// @Produce      json
// @Router       /subscriptions/notification [post]
// @Success      200  {object}  response.Common
func (c *SubscriptionController) HandlePaymentNotification(ctx *fiber.Ctx) error {
	// Get raw body
	body := ctx.Body()

	// Log the notification body for debugging
	utils.Log.Infof("Received payment notification: %s", string(body))

	// Try to extract order_id for debugging purposes
	var notificationData map[string]interface{}
	if err := json.Unmarshal(body, &notificationData); err != nil {
		utils.Log.Errorf("Failed to parse notification JSON: %v", err)
	} else {
		if orderID, ok := notificationData["order_id"].(string); ok {
			utils.Log.Infof("Notification for order ID: %s", orderID)
		} else {
			utils.Log.Warn("Notification does not contain order_id")
		}
	}

	// Process notification
	if err := c.Service.HandlePaymentNotification(ctx, body); err != nil {
		utils.Log.Errorf("Failed to process notification: %v", err)
		return utils.APIError(ctx, fiber.StatusInternalServerError, "notification_failed", err.Error())
	}

	// Return success response
	utils.Log.Info("Payment notification processed successfully")
	return ctx.JSON(response.Common{
		Status:  "success",
		Message: "Notification processed successfully",
	})
}
