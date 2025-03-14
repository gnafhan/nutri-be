package middleware

import (
	"app/src/config"
	"app/src/service"
	"app/src/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Auth(userService service.UserService, productTokenService service.ProductTokenService, requiredRights ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		productToken, err := productTokenService.GetProductTokenByUserID(c, user.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return fiber.ErrInternalServerError
		}

		if productToken == nil {
			return fiber.NewError(fiber.StatusForbidden, "Please activate your product token")
		}

		expDays := config.ProductTokenExpDays
		if expDays == "" {
			fmt.Println("Error: PRODUCT_TOKEN_EXP_DAYS not set")
			return fiber.ErrInternalServerError
		}

		expDuration, err := time.ParseDuration(expDays + "h")
		if err != nil {
			fmt.Println("Error parsing duration:", err)
			return fiber.ErrInternalServerError
		}

		expirationTime := productToken.ActivatedAt.Add(expDuration)
		if time.Now().After(expirationTime) {
			if err := productTokenService.DeleteProductToken(c, productToken.ID); err != nil {
				return fiber.ErrInternalServerError
			}
			return fiber.NewError(fiber.StatusForbidden, "Your product token has expired. Please activate a new one.")
		}

		c.Locals("user", user)

		if len(requiredRights) > 0 {
			userRights, hasRights := config.RoleRights[user.Role]
			if (!hasRights || !hasAllRights(userRights, requiredRights)) && c.Params("userId") != userID {
				return fiber.NewError(fiber.StatusForbidden, "You don't have permission to access this resource")
			}
		}

		return c.Next()
	}
}

func hasAllRights(userRights, requiredRights []string) bool {
	rightSet := make(map[string]struct{}, len(userRights))
	for _, right := range userRights {
		rightSet[right] = struct{}{}
	}

	for _, right := range requiredRights {
		if _, exists := rightSet[right]; !exists {
			return false
		}
	}
	return true
}
