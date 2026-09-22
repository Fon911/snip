package main

import (
	"fmt"
	"net/http"

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

	pool, err := storage.NewPool("postgres://postgres:secret@localhost:5400/postgres")

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
