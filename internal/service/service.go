package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"snap/example.com/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = storage.ErrNotFound

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(n int) string {
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		idx := rand.IntN(len(alphabet))
		result[i] = alphabet[idx]
	}
	return string(result)
}

func Shorten(ctx context.Context, pool *pgxpool.Pool, url string) (string, error) {
	var code string
	var err error
	for i := 0; i < 5; i++ {
		code = generateCode(8)
		err = storage.SaveLink(ctx, pool, code, url)
		if errors.Is(err, storage.ErrCodeTaken) {
			continue
		} else {
			break
		}
	}
	if err != nil {
		return "", err
	}

	return code, nil
}

func Stats(ctx context.Context, pool *pgxpool.Pool, code string) (string, int, error) {
	return storage.GetStats(ctx, pool, code)
}

func Resolve(ctx context.Context, pool *pgxpool.Pool, code string) (string, error) {
	url, err := storage.GetLink(ctx, pool, code)
	if err != nil {
		return url, err
	}
	err = storage.IncrementClicks(ctx, pool, code)
	if err != nil {
		fmt.Println(err)
	}
	return url, nil

}
