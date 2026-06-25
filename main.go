package main

import (
	"log"
	"net/http"
)

const (
	serverAddress = ":3000"
	baseURL       = "http://localhost:3000"
)

func main() {
	r := newRouter()

	log.Printf("Server started on %s", baseURL)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}
