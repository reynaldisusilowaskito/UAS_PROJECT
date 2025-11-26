package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	JWTSecret   string
	PostgresURI string
	MongoURI    string
}

var Env Config

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	Env = Config{
		AppPort:     os.Getenv("APP_PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		PostgresURI: os.Getenv("POSTGRES_URI"),
		MongoURI:    os.Getenv("MONGO_URI"),
	}
}
