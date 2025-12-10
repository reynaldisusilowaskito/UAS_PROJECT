package main

import (
	"log"

	"project_uas/config"
	"project_uas/database"
	"project_uas/route"

	"project_uas/app/repository"
	"project_uas/app/service"

	"github.com/gin-gonic/gin"
)

func main() {

	// Load environment variables
	config.LoadEnv()

	// Connect to PostgreSQL & MongoDB
	database.Connect()

	// Init Repository
	authRepo := repository.NewAuthRepo(database.PostgresDB)
	achievementRepo := repository.NewAchievementRepo(database.PostgresDB, database.MongoDB)
	studentRepo := repository.NewStudentRepo(database.PostgresDB)

	// Init Services
	authService := service.NewAuthService(authRepo)
	achievementService := service.NewAchievementService(achievementRepo, studentRepo)
	studentService := service.NewStudentService(studentRepo)

	// Init Router
	r := gin.Default()

	// Register Routes
	route.RegisterRoutes(r, authService, achievementService, studentService)

	
	for _, ri := range r.Routes() {
    log.Println(ri.Method, ri.Path)
}


	log.Println("Server running at :8080")
	r.Run(":8080")
}
