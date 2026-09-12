package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("нет такого кода")

func NewPool(dsn string) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("не смог создать пул: %w", err)
	}

	return pool, nil
}

func SaveLink(pool *pgxpool.Pool, code string, url string) error {
	_, err := pool.Exec(context.Background(),
		"INSERT INTO links (code, url) VALUES ($1, $2)",
		code, url)
	if err != nil {
		return fmt.Errorf("вставка упала: %w", err)
	}
	return nil
}

func GetLink(pool *pgxpool.Pool, code string) (string, error) {
	var url string
	err := pool.QueryRow(context.Background(),
		"SELECT url FROM links WHERE code = $1",
		code).Scan(&url)

	// Три исхода
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound

	}
	if err != nil {
		return "", fmt.Errorf("база сломалась: %w", err) // настоящий сбой

	}
	return url, nil
}

func CountLinks(pool *pgxpool.Pool) (int, error) {
	var count int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM links").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("не удалось посчитать запросы %w", err)
	}
	return count, nil
}
