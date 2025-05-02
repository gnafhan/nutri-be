package utils

import "github.com/gofiber/fiber/v2"

func APIError(c *fiber.Ctx, status int, errorCode string, message string, extras ...map[string]interface{}) error {
	response := fiber.Map{
		"status":  "error",
		"code":    errorCode,
		"message": message,
	}

	if len(extras) > 0 {
		for key, value := range extras[0] {
			response[key] = value
		}
	}

	return c.Status(status).JSON(response)
}
