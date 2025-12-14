package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"project_uas/helper"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid Authorization header",
			})
		}

		accessToken := parts[1]

		claims, err := helper.VerifyAccessToken(accessToken)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Set ke locals (Fiber)
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}
