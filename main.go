package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"project_uas/config"
	"project_uas/database"
	"project_uas/route"

	"project_uas/app/repository"
	"project_uas/app/service"

	_ "project_uas/docs"
)

func main() {

	// Load environment variables
	config.LoadEnv()

	// Connect to PostgreSQL & MongoDB
	database.Connect()

	// Init Repository
	authRepo := repository.NewAuthRepo(database.PostgresDB)
	achievementRepo := repository.NewAchievementRepo(
		database.PostgresDB,
		database.MongoDB,
	)
	studentRepo := repository.NewStudentRepo(database.PostgresDB)

	// Init Services
	authService := service.NewAuthService(authRepo)
	achievementService := service.NewAchievementService(
		achievementRepo,
		studentRepo,
	)
	studentService := service.NewStudentService(studentRepo)

	// Init Fiber App
	app := fiber.New()

	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Register Routes
	route.RegisterRoutes(
		app,
		authService,
		achievementService,
		studentService,
	)

	// Print registered routes
	for _, r := range app.GetRoutes() {
		log.Println(r.Method, r.Path)
	}

	log.Println("Server running at :3000")
	log.Fatal(app.Listen(":3000"))
}
