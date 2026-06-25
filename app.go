package main

import "github.com/jackc/pgx/v5/pgxpool"

type App struct {
	db *pgxpool.Pool
}
