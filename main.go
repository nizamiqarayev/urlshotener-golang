package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

const (
	serverAddress = ":3000"
	baseURL       = "http://localhost:3000"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ctx := context.Background()
	dbPool, err := connectDB(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbPool.Close()

	app := &App{db: dbPool}

	r := newRouter(app)

	log.Printf("Server started on %s", baseURL)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}
