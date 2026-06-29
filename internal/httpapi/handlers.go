package httpapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"urlShortener/internal/storage"
	"urlShortener/internal/urlutil"

	"github.com/go-chi/chi/v5"
	"github.com/lithammer/shortuuid/v4"
)

const maxShortKeyAttempts = 5

type createShortURLResponse struct {
	ShortURL    string `json:"short_url"`
	ShortKey    string `json:"short_key"`
	OriginalURL string `json:"original_url"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type StatsResponse struct {
	ShortKey    string    `json:"short_key"`
	OriginalURL string    `json:"original_url"`
	ClickCount  int64     `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := errorResponse{
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func (app *App) createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	originalURL := strings.TrimSpace(r.FormValue("url"))
	if originalURL == "" {
		writeErrorResponse(w, http.StatusBadRequest, "URL is required")
		return
	}

	originalURL, ok := urlutil.NormalizeHTTPURL(originalURL)
	if !ok {
		writeErrorResponse(w, http.StatusBadRequest, "URL must be a valid http or https URL with a real host")
		return
	}

	var key string
	inserted := false

	for attempt := 0; attempt < maxShortKeyAttempts; attempt++ {
		key = shortuuid.New()

		err := app.store.InsertMapping(r.Context(), key, originalURL)
		if err == nil {
			inserted = true
			break
		}

		if storage.IsUniqueViolation(err) {
			continue
		}

		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create short URL")
		return
	}

	if !inserted {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to create unique short URL")
		return
	}

	shortURL := fmt.Sprintf("%s/short/%s", app.baseURL, key)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := createShortURLResponse{
		ShortURL:    shortURL,
		ShortKey:    key,
		OriginalURL: originalURL,
	}

	log.Printf("Created short URL: %s -> %s", key, originalURL)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode create short URL response: %v", err)
	}
}

func (app *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Key is required")
		return
	}

	url, exists, err := app.store.FetchMapping(r.Context(), key)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch URL")
		return
	}

	if !exists {
		writeErrorResponse(w, http.StatusNotFound, "URL not found")
		return
	}

	if err := app.store.IncrementClickCount(r.Context(), key); err != nil {
		log.Printf("Failed to increment click count for %s: %v", key, err)
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status: "ok",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode health response: %v", err)
	}
}

func (app *App) readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := app.store.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		response := HealthResponse{
			Status: "not ready",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode readiness response: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := HealthResponse{
		Status: "ready",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode readiness response: %v", err)
	}
}

func (app *App) statsHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Key is required")
		return
	}

	stats, exists, err := app.store.GetStats(r.Context(), key)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch stats")
		return
	}

	if !exists {
		writeErrorResponse(w, http.StatusNotFound, "Stats not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := StatsResponse{
		ShortKey:    stats.ShortKey,
		OriginalURL: stats.OriginalURL,
		ClickCount:  stats.ClickCount,
		CreatedAt:   stats.CreatedAt,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode stats response: %v", err)
	}
}
