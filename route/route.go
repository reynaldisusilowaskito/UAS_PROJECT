package routes

import (
	"project_uas/database"
	"project_uas/app/repository"
	"project_uas/app/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	api := router.Group("/api/v1")

	
	// AUTHENTICATION
	
	auth := api.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/refresh", controller.RefreshToken)
		auth.POST("/logout", middleware.AuthMiddleware(), controller.Logout)
		auth.GET("/profile", middleware.AuthMiddleware(), controller.GetProfile)
	}

	// USERS (Admin Only)
	users := api.Group("/users", middleware.AuthMiddleware(), middleware.OnlyAdmin())
	{
		users.GET("/", controller.GetAllUsers)
		users.GET("/:id", controller.GetUserByID)
		users.POST("/", controller.CreateUser)
		users.PUT("/:id", controller.UpdateUser)
		users.DELETE("/:id", controller.DeleteUser)

		// Update Role User
		users.PUT("/:id/role", controller.UpdateUserRole)
	}

	//  ACHIEVEMENTS
	ach := api.Group("/achievements", middleware.AuthMiddleware())
	{
		ach.GET("/", controller.GetAchievements)             // list
		ach.GET("/:id", controller.GetAchievementDetail)     // detail

		// mahasiswa
		ach.POST("/", middleware.OnlyStudent(), controller.CreateAchievement)
		ach.PUT("/:id", middleware.OnlyStudent(), controller.UpdateAchievement)
		ach.DELETE("/:id", middleware.OnlyStudent(), controller.DeleteAchievement)

		// workflow: submit, verify, reject
		ach.POST("/:id/submit", middleware.OnlyStudent(), controller.SubmitAchievement)
		ach.POST("/:id/verify", middleware.OnlyLecturer(), controller.VerifyAchievement)
		ach.POST("/:id/reject", middleware.OnlyLecturer(), controller.RejectAchievement)

		// history
		ach.GET("/:id/history", controller.GetAchievementHistory)

		// file upload
		ach.POST("/:id/attachments", controller.UploadAttachment)
	}

	
	//  STUDENTS & LECTURERS

	students := api.Group("/students", middleware.AuthMiddleware())
	{
		students.GET("/", controller.GetAllStudents)
		students.GET("/:id", controller.GetStudentByID)
		students.GET("/:id/achievements", controller.GetStudentAchievements)
		students.PUT("/:id/advisor", middleware.OnlyAdmin(), controller.UpdateStudentAdvisor)
	}

	lecturers := api.Group("/lecturers", middleware.AuthMiddleware())
	{
		lecturers.GET("/", controller.GetAllLecturers)
		lecturers.GET("/:id/advisees", controller.GetLecturerAdvisees)
	}

	//  REPORTS & ANALYTICS
	reports := api.Group("/reports", middleware.AuthMiddleware())
	{
		reports.GET("/statistics", middleware.OnlyAdmin(), controller.GetStatisticsReport)
		reports.GET("/student/:id", controller.GetStudentReport)
	}
}
