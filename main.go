package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mohsenjafari-aiio/aiiobackend/internal/config"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	// Connect to database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get underlying SQL DB to close connection when done
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	log.Println("Application started successfully!")
	log.Println("Database connection established")

	// TODO: Add your application logic here
	// For example: start HTTP server, initialize repositories, etc.
}
