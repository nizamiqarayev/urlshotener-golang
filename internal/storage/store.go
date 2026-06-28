package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (store *Store) InsertMapping(ctx context.Context, key, originalURL string) error {
	_, err := store.db.Exec(ctx, `
		INSERT INTO urls (short_key, original_url)
		VALUES ($1, $2)
	`, key, originalURL)

	return err
}

func (store *Store) FetchMapping(ctx context.Context, key string) (string, bool, error) {
	var originalURL string

	err := store.db.QueryRow(ctx, `
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

func (store *Store) IncrementClickCount(ctx context.Context, key string) error {
	_, err := store.db.Exec(ctx, `
		UPDATE urls
		SET click_count = click_count + 1
		WHERE short_key = $1
	`, key)

	return err
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
