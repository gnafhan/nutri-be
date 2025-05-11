package middleware

import (
	"app/src/model"
	"app/src/service"
	"app/src/utils"

	"github.com/gofiber/fiber/v2"
)

func SubscriptionRequired(subService service.SubscriptionService, feature_name string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*model.User)

		_, err := subService.GetUserActiveSubscription(c, user.ID)
		if err != nil {
			return utils.APIError(c, fiber.StatusForbidden,
				"subscription_required",
				"Active subscription required",
				map[string]interface{}{
					"upgrade_url": "/v1/subscriptions/plans",
				})
		}
		if feature_name != "" {
			// Check if the user has access to the requested feature
			hasAccess, err := subService.CheckFeatureAccess(c, user.ID, feature_name)
			if err != nil {
				return utils.APIError(c, fiber.StatusForbidden,
					"subscription_required",
					"Active subscription required",
					map[string]interface{}{
						"upgrade_url": "/v1/subscriptions/plans",
					})
			}
			if !hasAccess {
				return utils.APIError(c, fiber.StatusForbidden,
					"access",
					"You don't have access to this feature",
					map[string]interface{}{
						"upgrade_url": "/v1/subscriptions/plans",
					})
			}
		}

		return c.Next()
	}
}
