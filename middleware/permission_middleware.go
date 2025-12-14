package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		perms, ok := c.Locals("permissions").([]string)
		if !ok {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "permissions not found",
			})
		}

		for _, p := range perms {
			if p == permission {
				return c.Next()
			}
		}

		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "missing permission: " + permission,
		})
	}
}
