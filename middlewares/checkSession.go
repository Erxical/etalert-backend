package middlewares

import (
	"etalert-backend/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ValidateSession(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the access token from the Authorization header
		accessToken := c.Get("Authorization")
		if accessToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing access token"})
		}

		parts := strings.Split(accessToken, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid authorization header format"})
		}

		// Validate the access token
		claims, err := authService.ValidateAccessToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired access token"})
		}

		// Store the user ID in the context for further use in handlers
		c.Locals("userId", claims["googleId"])

		// Proceed to the next middleware/handler
		return c.Next()
	}
}
