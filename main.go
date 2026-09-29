package main

import (
	"fmt"
	"net/http"
	"os"

	"snap/example.com/internal/handler"
	"snap/example.com/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {

	rt := chi.NewRouter()
	rt.Get("/{code}", handler.Get)
	rt.Get("/api/stats/{code}", handler.Stats)
	rt.Post("/api/shorten", handler.Shorten)

	// 1. Связь с базой

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:secret@localhost:5400/postgres"
	}
	pool, err := storage.NewPool(dsn)

	if err != nil {
		fmt.Println(err)
		return
	}
	handler.Pool = pool
	defer pool.Close()
	fmt.Println("база подключена")

	// 3. Ручки

	// 4. Дежурство — ВСЕГДА последняя строка
	err = http.ListenAndServe(":8000", rt)
	if err != nil {
		fmt.Println(err)
	}
}
