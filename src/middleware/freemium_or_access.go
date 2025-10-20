package middleware

import (
	"app/src/config"
	"app/src/service"
	"app/src/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FreemiumOrAccess(userService service.UserService, productTokenService service.ProductTokenService, subscriptionService service.SubscriptionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Verify JWT token (required)
		authHeader := c.Get("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
		}

		userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
		}

		user, err := userService.GetUserByID(c, userID)
		if err != nil || user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
		}

		// 2. Check if user has active subscription (including freemium) OR valid product token
		hasAccess := false
		accessType := ""

		// Check for active subscription first (including freemium)
		subscription, err := subscriptionService.GetUserActiveSubscription(c, user.ID)
		if err == nil && subscription != nil {
			// Check if subscription is still active (not expired)
			if subscription.IsActive && time.Now().Before(subscription.EndDate) {
				hasAccess = true
				accessType = "subscription"

				// Check if it's a freemium subscription
				if subscription.Plan.Name == "Freemium Trial" {
					accessType = "freemium"
				}
			}
		}

		// If no active subscription, check for product token
		if !hasAccess {
			productToken, err := productTokenService.GetProductTokenByUserID(c, user.ID)
			if err == nil && productToken != nil {
				// Check if product token is still valid (not expired)
				expDays := config.ProductTokenExpDays
				if expDays != "" {
					if expDaysInt, parseErr := strconv.Atoi(expDays); parseErr == nil {
						expDuration := time.Duration(expDaysInt) * 24 * time.Hour
						expirationTime := productToken.ActivatedAt.Add(expDuration)
						if time.Now().Before(expirationTime) {
							hasAccess = true
							accessType = "product_token"
						}
					}
				}
			}
		}

		// 3. If neither subscription nor product token: return 403
		if !hasAccess {
			return utils.APIError(c, fiber.StatusForbidden,
				"access_required",
				"Please activate your product token or subscribe to access this feature",
				map[string]interface{}{
					"upgrade_url": "/v1/subscriptions/plans",
				})
		}

		// 4. If freemium expired: return 403 with upgrade message
		if accessType == "freemium" && subscription != nil && time.Now().After(subscription.EndDate) {
			return utils.APIError(c, fiber.StatusForbidden,
				"freemium_expired",
				"Your free trial has expired. Please upgrade to continue using the service",
				map[string]interface{}{
					"upgrade_url": "/v1/subscriptions/plans",
					"expired_at":  subscription.EndDate.Format(time.RFC3339),
				})
		}

		// 5. Store user in context and proceed
		c.Locals("user", user)
		c.Locals("access_type", accessType)

		return c.Next()
	}
}
