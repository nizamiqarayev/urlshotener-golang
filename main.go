package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lithammer/shortuuid/v4"
)

const (
	serverAddress = ":3000"
	baseURL       = "http://localhost:3000"
)

type Mapper struct {
	mapping map[string]string
	lock    sync.Mutex
}

var urlMapper = newMapper()

func main() {
	r := newRouter()

	log.Printf("Server started on %s", baseURL)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}

func newRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.Post("/shorten", createShortURLHandler)
	r.Get("/short/{key}", redirectHandler)

	return r
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

func newMapper() Mapper {
	return Mapper{
		mapping: make(map[string]string),
		lock:    sync.Mutex{},
	}
}

func insertMapping(key, originalURL string) {
	urlMapper.lock.Lock()
	defer urlMapper.lock.Unlock()

	urlMapper.mapping[key] = originalURL
}

func fetchMapping(key string) (string, bool) {
	urlMapper.lock.Lock()
	defer urlMapper.lock.Unlock()

	originalURL, exists := urlMapper.mapping[key]
	return originalURL, exists
}
