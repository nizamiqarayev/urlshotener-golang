package httpapi

import "urlShortener/internal/storage"

type App struct {
	baseURL string
	store   *storage.Store
}
