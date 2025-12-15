package route

import (
	"github.com/gofiber/fiber/v2"

	"project_uas/app/service"
	"project_uas/middleware"
)

func RegisterRoutes(
	app *fiber.App,
	authService *service.AuthService,
	achievementService *service.AchievementService,
	studentService *service.StudentService,
) {

	api := app.Group("/api/v1")

	// =====================
	// AUTH
	// =====================
	auth := api.Group("/auth")
	{
		auth.Post("/login", authService.Login)
		auth.Post("/refresh", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotImplemented) })
		auth.Post("/logout", middleware.AuthMiddleware(), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})
		auth.Get("/profile", middleware.AuthMiddleware(), authService.GetProfile)
	}

	// =====================
	// USERS
	// =====================
	users := api.Group("/users", middleware.AuthMiddleware())
	{
		users.Get("/", middleware.RequirePermission("users:read"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNotImplemented)
		})
		users.Get("/:id", middleware.RequirePermission("users:read"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNotImplemented)
		})
		users.Post("/", middleware.RequirePermission("users:create"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNotImplemented)
		})
		users.Put("/:id", middleware.RequirePermission("users:update"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNotImplemented)
		})
		users.Delete("/:id", middleware.RequirePermission("users:delete"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNotImplemented)
		})
		users.Put("/:id/role", middleware.RequirePermission("users:update-role"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNotImplemented)
		})
	}

	// =====================
	// ACHIEVEMENTS
	// =====================
	ach := api.Group("/achievements", middleware.AuthMiddleware())
	{
		ach.Post("/", middleware.OnlyStudent(), achievementService.CreateAchievement)
		ach.Post("/:id/submit", middleware.OnlyStudent(), achievementService.SubmitAchievement)
		ach.Post("/:id/verify", middleware.OnlyLecturer(), achievementService.VerifyAchievement)
		ach.Post("/:id/reject", middleware.OnlyLecturer(), achievementService.RejectAchievement)
		ach.Get("/:id/history", achievementService.GetAchievementHistory)
		ach.Post("/:id/attachments", achievementService.UploadAttachment)
		ach.Get("/:id", achievementService.GetAchievementDetail)
		ach.Delete("/:id", middleware.OnlyStudent(), achievementService.DeleteAchievement)
	}

	// =====================
	// STUDENTS
	// =====================
	students := api.Group("/students", middleware.AuthMiddleware())
	{
		students.Get("/", func(c *fiber.Ctx) error {return c.SendStatus(fiber.StatusNotImplemented)})
		students.Get("/:id", func(c *fiber.Ctx) error {return c.SendStatus(fiber.StatusNotImplemented)})
		students.Get("/:id/achievements", func(c *fiber.Ctx) error {return c.SendStatus(fiber.StatusNotImplemented)})
		students.Put("/:id/advisor", middleware.OnlyAdmin(), func(c *fiber.Ctx) error {return c.SendStatus(fiber.StatusNotImplemented)})
		students.Get("/profile", studentService.GetProfile)
	}

	// =====================
	// LECTURERS
	// =====================
	lecturers := api.Group("/lecturers", middleware.AuthMiddleware())
	{
		lecturers.Get("/", func(c *fiber.Ctx) error {return c.SendStatus(fiber.StatusNotImplemented)})
		lecturers.Get("/:id/advisees",middleware.OnlyLecturer(),achievementService.GetAdviseeAchievements,)}

	// =====================
	// REPORTS
	// =====================
	reports := api.Group("/reports", middleware.AuthMiddleware())
	{
		reports.Get("/statistics", middleware.OnlyAdmin(), func(c *fiber.Ctx) error {return c.SendStatus(fiber.StatusNotImplemented)})
		reports.Get("/student/:id", func(c *fiber.Ctx) error {return c.SendStatus(fiber.StatusNotImplemented)})
	}
}
