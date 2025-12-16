package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"project_uas/helper"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "missing Authorization header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid Authorization header"})
		}

		token := parts[1]

		// ✅ STEP 1: CEK BLACKLIST (via helper)
		if helper.IsTokenBlacklisted(token) {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "token has been revoked"})
		}

		// ✅ STEP 2: CEK JWT
		claims, err := helper.VerifyAccessToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid or expired token"})
		}

		// ✅ SET CONTEXT
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}
