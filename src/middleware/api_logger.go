package middleware

import (
	"app/src/model"
	"app/src/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// APILoggerConfig creates middleware that logs API requests and responses
func APILoggerConfig() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate a unique request ID
		requestID := fmt.Sprintf("req-%s", uuid.New().String()[:8])
		c.Locals("requestID", requestID)

		// Start timer
		start := time.Now()

		// Get request info
		method := c.Method()
		path := c.Path()
		ip := c.IP()
		userAgent := c.Get("User-Agent")

		// Parse request body without consuming it
		var requestBody interface{}
		if method == "POST" || method == "PUT" || method == "PATCH" {
			body := c.Body()
			if len(body) > 0 {
				// Try to parse as JSON
				if err := json.Unmarshal(body, &requestBody); err != nil {
					requestBody = string(body)
				}
				// Reset body for downstream handlers
				c.Request().SetBody(body)
			}
		}

		// Get user ID if authenticated
		var userID string
		user := c.Locals("user")
		if user != nil {
			if u, ok := user.(*model.User); ok && u != nil {
				userID = u.ID.String()
			}
		}

		// Log request
		utils.LogAPIRequest(utils.RequestResponseData{
			Method:      method,
			Path:        path,
			RequestID:   requestID,
			IPAddress:   ip,
			UserAgent:   userAgent,
			RequestBody: requestBody,
			UserID:      userID,
		})

		// Get the original response body
		originalResponseBody := c.Response().Body()

		// Process the request
		err := c.Next()

		// Calculate request duration
		duration := time.Since(start)
		elapsedTime := fmt.Sprintf("%s", duration)

		// Get response status code
		statusCode := c.Response().StatusCode()

		// Get response body
		var responseBody interface{}
		if len(originalResponseBody) > 0 {
			// Try to parse as JSON
			if err := json.Unmarshal(originalResponseBody, &responseBody); err != nil {
				responseBody = string(originalResponseBody)
			}
		}

		// Log response
		utils.LogAPIResponse(utils.RequestResponseData{
			Method:      method,
			Path:        path,
			RequestID:   requestID,
			StatusCode:  statusCode,
			Response:    responseBody,
			ElapsedTime: elapsedTime,
			UserID:      userID,
		})

		// Log user activity if authenticated
		if userID != "" {
			utils.LogUserActivity(utils.ActivityData{
				UserID:      userID,
				Action:      method,
				Resource:    path,
				RequestID:   requestID,
				IPAddress:   ip,
				UserAgent:   userAgent,
				StatusCode:  statusCode,
				ElapsedTime: elapsedTime,
			})
		}

		return err
	}
}
