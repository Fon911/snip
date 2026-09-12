package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprintln(w, "Hello input name  /Hello?name=Lexa", name)
		return
	}
	fmt.Fprintln(w, "Hello world ", name)
}

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
	key := r.URL.Query().Get("code")
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
	http.HandleFunc("/Hello", Hello)
	http.HandleFunc("/shorten", Shorten)
	http.HandleFunc("/go", get)

	// 4. Дежурство — ВСЕГДА последняя строка
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
