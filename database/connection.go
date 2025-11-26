package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"project_uas/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global Vars
var (
	PostgresDB *sqlx.DB
	MongoDB    *mongo.Database
)

// MAIN CONNECT FUNCTION (dipanggil di main.go)
func Connect() {
	config.LoadEnv()

	connectPostgres()
	connectMongo()
}

// =====================================================
// ===============   POSTGRESQL CONNECT   ===============
// =====================================================
func connectPostgres() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("FAILED CONNECT POSTGRES:", err)
	}

	PostgresDB = db
	log.Println("CONNECTED TO POSTGRES!")
}

// =====================================================
// ==================   MONGO CONNECT   =================
// =====================================================
func connectMongo() {
	mongoURI := fmt.Sprintf(
		"mongodb://%s:%s/%s",
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DB"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("FAILED CONNECT MONGODB:", err)
	}

	// Ping
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("FAILED PING MONGODB:", err)
	}

	MongoDB = client.Database(os.Getenv("MONGO_DB"))
	log.Println("CONNECTED TO MONGODB!")
}
