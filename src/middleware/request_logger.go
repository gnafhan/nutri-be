package middleware

import (
	"app/src/utils"

	"github.com/gofiber/fiber/v2"
)

// RequestBodyLoggerConfig creates middleware that logs request bodies
func RequestBodyLoggerConfig() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only log request body for specific methods that typically have bodies
		method := c.Method()
		if method == "POST" || method == "PUT" || method == "PATCH" {
			// Get request path
			path := c.Path()

			// Get client IP
			ip := c.IP()

			// Clone request body so we can read it without consuming it
			body := c.Body()
			if len(body) > 0 {
				// Parse body as JSON
				var bodyMap map[string]interface{}
				if err := c.BodyParser(&bodyMap); err != nil {
					// If body parsing fails, log the raw body
					utils.LogRequestBody(path, method, string(body), ip)
				} else {
					// Log the parsed JSON body
					utils.LogRequestBody(path, method, bodyMap, ip)
				}

				// Re-set the body for downstream handlers
				c.Request().SetBody(body)
			}
		}

		return c.Next()
	}
}
