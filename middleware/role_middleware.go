package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func OnlyAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("role") != "admin" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "admin only",
			})
		}
		return c.Next()
	}
}

func OnlyStudent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("role") != "student" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "students only",
			})
		}
		return c.Next()
	}
}

func OnlyLecturer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("role") != "lecturer" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "lecturers only",
			})
		}
		return c.Next()
	}
}

// ================================
//   SELF ACCESS VALIDATION
// ================================
func OnlySelf() fiber.Handler {
	return func(c *fiber.Ctx) error {

		role := c.Locals("role")
		userID := c.Locals("user_id")
		paramID := c.Params("id")

		if role == "admin" {
			return c.Next()
		}

		if userID != paramID {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "you can only access your own account",
			})
		}

		return c.Next()
	}
}

func OnlyStudentSelf() fiber.Handler {
	return func(c *fiber.Ctx) error {

		if c.Locals("role") != "student" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "students only",
			})
		}

		if c.Locals("user_id") != c.Params("id") {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "you can only access your own account",
			})
		}

		return c.Next()
	}
}

func OnlyLecturerSelf() fiber.Handler {
	return func(c *fiber.Ctx) error {

		if c.Locals("role") != "lecturer" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "lecturers only",
			})
		}

		if c.Locals("user_id") != c.Params("id") {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "you can only access your own account",
			})
		}

		return c.Next()
	}
}
