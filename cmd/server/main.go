package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"urlShortener/internal/database"
	"urlShortener/internal/httpapi"

	"github.com/joho/godotenv"
)

var serverAddress = ":3000"

var baseURL = "http://localhost:3000"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if value := os.Getenv("BASE_URL"); value != "" {
		baseURL = value
	}
	if value := os.Getenv("SERVER_ADDRESS"); value != "" {
		serverAddress = value
	}

	ctx := context.Background()
	dbPool, err := database.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbPool.Close()

	router := httpapi.NewRouter(dbPool, baseURL)

	log.Printf("Server started on %s", baseURL)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
