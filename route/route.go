package route

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authService *service.AuthService, achievementService *service.AchievementService, studentService *service.StudentService) {

	api := router.Group("/api/v1")

	// =====================
	// AUTH
	// =====================
	auth := api.Group("/auth")
	{
		auth.POST("/login", authService.Login)
		auth.POST("/refresh", nil)
		auth.POST("/logout", middleware.AuthMiddleware(), nil)
		auth.GET("/profile", middleware.AuthMiddleware(), authService.GetProfile)
	}

	// =====================
	// USERS
	// =====================
	users := api.Group("/users", middleware.AuthMiddleware())
	{
		users.GET("/", middleware.RequirePermission("users:read"), nil)
		users.GET("/:id", middleware.RequirePermission("users:read"), nil)
		users.POST("/", middleware.RequirePermission("users:create"), nil)
		users.PUT("/:id", middleware.RequirePermission("users:update"), nil)
		users.DELETE("/:id", middleware.RequirePermission("users:delete"), nil)
		users.PUT("/:id/role", middleware.RequirePermission("users:update-role"), nil)
	}

	// =====================
	// ACHIEVEMENTS
	// =====================
	ach := api.Group("/achievements", middleware.AuthMiddleware())
	{
		ach.POST("/", middleware.OnlyStudent(), achievementService.CreateAchievement)
		ach.POST("/:id/submit", middleware.OnlyStudent(), achievementService.SubmitAchievement)
		ach.POST("/:id/verify", middleware.OnlyLecturer(), achievementService.VerifyAchievement)
		ach.POST("/:id/reject", middleware.OnlyLecturer(), achievementService.RejectAchievement)
		ach.GET("/:id/history", achievementService.GetAchievementHistory)
		ach.POST("/:id/attachments", achievementService.UploadAttachment)
		ach.GET("/:id", achievementService.GetAchievementDetail)
		ach.DELETE("/:id", middleware.OnlyStudent(), achievementService.DeleteAchievement)

	}

	// =====================
	// STUDENTS
	// =====================
	students := api.Group("/students", middleware.AuthMiddleware())
	{
		students.GET("/", nil)
		students.GET("/:id", nil)
		students.GET("/:id/achievements", nil)
		students.PUT("/:id/advisor", middleware.OnlyAdmin(), nil)
		students.GET("/profile", studentService.GetProfile)
	}

	// =====================
	// LECTURERS
	// =====================
	lecturers := api.Group("/lecturers", middleware.AuthMiddleware())
	{
		lecturers.GET("/", nil)
		lecturers.GET("/:id/advisees", nil)
	}

	// =====================
	// REPORTS
	// =====================
	reports := api.Group("/reports", middleware.AuthMiddleware())
	{
		reports.GET("/statistics", middleware.OnlyAdmin(), nil)
		reports.GET("/student/:id", nil)
	}
}
