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
	userService *service.UserService,
	lecturerService *service.LecturerService,
	reportService *service.ReportService,
) {

	api := app.Group("/api/v1")

	// =====================
	// AUTH
	// =====================
	auth := api.Group("/auth")
	{
		auth.Post("/login", authService.Login)
		auth.Post("/refresh", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotImplemented) })
		auth.Post("/logout",middleware.AuthMiddleware(),authService.Logout,)
		auth.Get("/profile", middleware.AuthMiddleware(), authService.GetProfile)
	}

	// =====================
	// USERS
	// =====================
	users := api.Group("/users", middleware.AuthMiddleware())
{
	users.Get("/", middleware.RequirePermission("users:read"), userService.GetAll)
	users.Get("/:id", middleware.RequirePermission("users:read"), userService.GetByID)
	users.Post("/", middleware.RequirePermission("users:create"), userService.Create)
	users.Put("/:id", middleware.RequirePermission("users:update"), userService.Update)
	users.Delete("/:id", middleware.RequirePermission("users:delete"), userService.Delete)
	users.Put("/:id/role", middleware.RequirePermission("users:update-role"), userService.UpdateRole)
}


	// =====================
	// ACHIEVEMENTS
	// =====================
	ach := api.Group("/achievements", middleware.AuthMiddleware())
	{
		ach.Get("/", middleware.OnlyAdmin(), achievementService.GetAll)
		ach.Post("/", middleware.OnlyStudent(), achievementService.CreateAchievement)
		ach.Post("/:id/submit", middleware.OnlyStudent(), achievementService.SubmitAchievement)
		ach.Post("/:id/verify", middleware.OnlyLecturer(), achievementService.VerifyAchievement)
		ach.Post("/:id/reject", middleware.OnlyLecturer(), achievementService.RejectAchievement)
		ach.Get("/:id/history", achievementService.GetAchievementHistory)
		ach.Post("/:id/attachments", achievementService.UploadAttachment)
		ach.Get("/:id", achievementService.GetAchievementDetail)
		ach.Delete("/:id", middleware.OnlyStudent(), achievementService.DeleteAchievement)
		ach.Put("/:id", middleware.OnlyStudent(), achievementService.UpdateAchievement)
	}

	// =====================
	// STUDENTS
	// =====================
	students := api.Group("/students", middleware.AuthMiddleware())
{
		students.Get("/", studentService.GetAll)
		students.Get("/:id", studentService.GetByID)
		students.Get("/:id/achievements", studentService.GetAchievements)
		students.Put("/:id/advisor", middleware.OnlyAdmin(), studentService.UpdateAdvisor)
		students.Get("/profile", studentService.GetProfile)
}


	// =====================
	// LECTURERS
	// =====================
	lecturers := api.Group("/lecturers", middleware.AuthMiddleware())
{
	lecturers.Get("/", lecturerService.GetAll)
	lecturers.Get("/profile", middleware.OnlyLecturer(), lecturerService.GetProfile)
	lecturers.Get("/:id/advisees",middleware.OnlyLecturer(),achievementService.GetAdviseeAchievements,)
}

	// =====================
	// REPORTS
	// =====================
	reports := api.Group("/reports", middleware.AuthMiddleware())
{
	reports.Get("/student/:id",middleware.OnlyAdmin(),reportService.GetStudentReport,)
	reports.Get("/statistics",middleware.OnlyAdmin(),reportService.GetAchievementStats,)
}

}
