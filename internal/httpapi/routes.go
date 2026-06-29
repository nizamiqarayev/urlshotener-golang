package httpapi

import (
	"net/http"

	"urlShortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, baseURL string) http.Handler {
	app := &App{
		baseURL: baseURL,
		store:   storage.NewStore(db),
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/health", app.healthHandler)
	r.Get("/healthz", app.healthHandler)
	r.Post("/shorten", app.createShortURLHandler)
	r.Get("/short/{key}", app.redirectHandler)
	r.Get("/ready", app.readyHandler)
	r.Get("/stats/{key}", app.statsHandler)

	return r
}
