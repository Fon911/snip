package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

var KeyCount int

func Shorten(w http.ResponseWriter, r *http.Request) {

	val := r.URL.Query().Get("url")
	if val == "" {
		fmt.Fprintln(w, "нужно так: /shorten?url=ссылка")
		return
	}
	KeyCount++
	key := strconv.Itoa(KeyCount)
	err := SaveLink(pool, key, val)
	if err != nil {
		fmt.Fprintln(w, "не смог сохранить")
		return
	}
	fmt.Fprintln(w, "записал:", key)
}

func get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "code")
	url, err := GetLink(pool, key)
	if errors.Is(err, ErrNotFound) {
		fmt.Fprintln(w, "нет такого кода")
		return
	}
	if err != nil {
		fmt.Fprintln(w, "Ошибка сервера")
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {

	rt := chi.NewRouter()
	rt.Get("/{code}", get)

	rt.Post("/api/shorten", Shorten)

	// 1. Связь с базой
	p, err := NewPool("postgres://postgres:secret@localhost:5400/postgres")
	if err != nil {
		fmt.Println(err)
		return
	}
	pool = p
	defer pool.Close()
	fmt.Println("база подключена")

	// 2. Счётчик
	CountLink, ErrCountLink := CountLinks(pool)
	if ErrCountLink != nil {
		fmt.Println(ErrCountLink)
		return
	}
	KeyCount = CountLink
	fmt.Println("ссылок в базе:", KeyCount)

	// 3. Ручки

	// 4. Дежурство — ВСЕГДА последняя строка
	err = http.ListenAndServe(":8000", rt)
	if err != nil {
		fmt.Println(err)
	}
}
