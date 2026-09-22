package storage

import (
	"context"
	"errors"
	"fmt"
	"snap/example.com/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("нет такого кода")
var ErrCodeTaken = errors.New("код занят")

func NewPool(dsn string) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("не смог создать пул: %w", err)
	}

	return pool, nil
}

func SaveLink(pool *pgxpool.Pool, code string, url string) error {
	q := db.New(pool)
	err := q.SaveLink(context.Background(), db.SaveLinkParams{
		Code: code,
		Url:  url,
	})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrCodeTaken
	}

	if err != nil {
		return fmt.Errorf("вставка упала: %w", err)
	}
	return nil
}

func GetLink(pool *pgxpool.Pool, code string) (string, error) {

	q := db.New(pool)
	url, err := q.GetLink(context.Background(), code)
	// Три исхода
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("база сломалась: %w", err) // настоящий сбой

	}

	return url, nil
}

func IncrementClicks(pool *pgxpool.Pool, code string) error {
	q := db.New(pool)
	err := q.IncrementClicks(context.Background(), code)
	if err != nil {
		return fmt.Errorf("не смог посчитать клик: %w", err)
	}
	return nil

}

func GetStats(pool *pgxpool.Pool, code string) (string, int, error) {

	q := db.New(pool)
	row, err := q.GetStats(context.Background(), code)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", 0, ErrNotFound
	}
	if err != nil {
		return "", 0, fmt.Errorf("база сломалась: %w", err)
	}
	return row.Url, int(row.Clicks), nil
}
