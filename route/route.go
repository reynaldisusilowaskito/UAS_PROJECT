package routes

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authService *service.AuthService) {

	api := router.Group("/api/v1")

	// =====================
	// AUTH
	// =====================
	auth := api.Group("/auth")
	{
		auth.POST("/login", authService.Login)
		auth.POST("/refresh", nil)
		auth.POST("/logout", middleware.AuthMiddleware(), nil)
		auth.GET("/profile", middleware.AuthMiddleware(), nil)
	}

	// =====================
	// USERS
	// =====================
	users := api.Group("/users", middleware.AuthMiddleware(), middleware.OnlyAdmin())
	{
		users.GET("/", nil)
		users.GET("/:id", nil)
		users.POST("/", nil)
		users.PUT("/:id", nil)
		users.DELETE("/:id", nil)
		users.PUT("/:id/role", nil)
	}

	// =====================
	// ACHIEVEMENTS
	// =====================
	ach := api.Group("/achievements", middleware.AuthMiddleware())
	{
		ach.GET("/", nil)
		ach.GET("/:id", nil)

		ach.POST("/", middleware.OnlyStudent(), nil)
		ach.PUT("/:id", middleware.OnlyStudent(), nil)
		ach.DELETE("/:id", middleware.OnlyStudent(), nil)

		ach.POST("/:id/submit", middleware.OnlyStudent(), nil)
		ach.POST("/:id/verify", middleware.OnlyLecturer(), nil)
		ach.POST("/:id/reject", middleware.OnlyLecturer(), nil)

		ach.GET("/:id/history", nil)
		ach.POST("/:id/attachments", nil)
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
