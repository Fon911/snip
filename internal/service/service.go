package service

import (
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

func Shorten(pool *pgxpool.Pool, url string) (string, error) {
	var code string
	var err error
	for i := 0; i < 5; i++ {
		code = generateCode(8)
		err = storage.SaveLink(pool, code, url)
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

func Stats(pool *pgxpool.Pool, code string) (string, int, error) {
	return storage.GetStats(pool, code)
}

func Resolve(pool *pgxpool.Pool, code string) (string, error) {
	url, err := storage.GetLink(pool, code)
	if err != nil {
		return url, err
	}
	err = storage.IncrementClicks(pool, code)
	if err != nil {
		fmt.Println(err)
	}
	return url, nil

}
