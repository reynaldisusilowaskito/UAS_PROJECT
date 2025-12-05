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

	// Run migration
	// database.Migrate(database.PostgresDB)

	// Init Repository
	authRepo := repository.NewAuthRepo(database.PostgresDB)

	// Init Service
	authService := service.NewAuthService(authRepo)

	// Init Router
	r := gin.Default()

	// Register Routes (inject service)
	routes.RegisterRoutes(r, authService)

	log.Println("Server running at :8080")
	r.Run(":8080")
}
