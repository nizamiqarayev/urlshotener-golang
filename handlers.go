package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lithammer/shortuuid/v4"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	key := shortuuid.New()
	shortURL := fmt.Sprintf("%s/short/%s", baseURL, key)

	insertMapping(key, originalURL)
	w.WriteHeader(http.StatusCreated)
	log.Printf("Created short URL: %s -> %s", key, originalURL)
	w.Write([]byte("Created short URL: " + shortURL))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	url, exists := fetchMapping(key)
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
