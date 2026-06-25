package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newRouter(app *App) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Post("/shorten", app.createShortURLHandler)
	r.Get("/short/{key}", app.redirectHandler)

	return r
}
