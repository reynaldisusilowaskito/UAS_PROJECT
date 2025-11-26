package main

import (
	"project_uas/config"
	"project_uas/database"
)

func main() {
	// load .env
	config.LoadEnv()

	// connect to PostgreSQL

	// connect to MongoDB
	database.Connect()
}
