package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insertMapping(ctx context.Context, db *pgxpool.Pool, key, originalURL string) error {
	_, err := db.Exec(ctx, `
		INSERT INTO urls (short_key, original_url)
		VALUES ($1, $2)
	`, key, originalURL)

	return err
}

func fetchMapping(ctx context.Context, db *pgxpool.Pool, key string) (string, bool, error) {
	var originalURL string

	err := db.QueryRow(ctx, `
		SELECT original_url
		FROM urls
		WHERE short_key = $1
	`, key).Scan(&originalURL)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}

		return "", false, err
	}

	return originalURL, true, nil
}

func incrementClickCount(ctx context.Context, db *pgxpool.Pool, key string) error {
	_, err := db.Exec(ctx, `
		UPDATE urls
		SET click_count = click_count + 1
		WHERE short_key = $1
	`, key)

	return err
}
