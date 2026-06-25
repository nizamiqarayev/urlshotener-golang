package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Post("/shorten", createShortURLHandler)
	r.Get("/short/{key}", redirectHandler)

	return r
}
