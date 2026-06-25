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

func (app *App) createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	key := shortuuid.New()
	shortURL := fmt.Sprintf("%s/short/%s", baseURL, key)

	err := insertMapping(r.Context(), app.db, key, originalURL)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("Created short URL: %s -> %s", key, originalURL)
	w.Write([]byte("Created short URL: " + shortURL))
}

func (app *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	url, exists, err := fetchMapping(r.Context(), app.db, key)
	if err != nil {
		http.Error(w, "Failed to fetch URL", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	if err := incrementClickCount(r.Context(), app.db, key); err != nil {
		log.Printf("Failed to increment click count for %s: %v", key, err)
	}

	http.Redirect(w, r, url, http.StatusFound)
}
