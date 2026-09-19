package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Code string `json:"code"`
}

type StatResponse struct {
	Code   string `json:"code"`
	URL    string `json:"url"`
	Clicks int    `json:"clicks"`
}

func Shorten(w http.ResponseWriter, r *http.Request) {

	var req ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "кривой JSON")
		return
	}
	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "нужно поле URL")
		return
	}
	var code string
	for i := 0; i < 5; i++ {
		code = GenerateCode(8)
		err = SaveLink(pool, code, req.URL)
		if errors.Is(err, ErrCodeTaken) {
			continue
		} else {
			break
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "не смог сохранить")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ShortenResponse{Code: code})
}

func get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	url, err := GetLink(pool, code)
	if errors.Is(err, ErrNotFound) {
		fmt.Fprintln(w, "нет такого кода")
		return
	}
	if err != nil {
		fmt.Fprintln(w, "Ошибка сервера")
		return
	}
	err = IncrementClicks(pool, code)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, url, http.StatusFound)

}

func Stats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	url, clicks, err := GetStats(pool, code)
	if errors.Is(err, ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "нет такоq ссылки")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "ошибка сервера")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StatResponse{Code: code, URL: url, Clicks: clicks})
}

func main() {

	rt := chi.NewRouter()
	rt.Get("/{code}", get)
	rt.Get("/api/stats/{code}", Stats)
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

	// 3. Ручки

	// 4. Дежурство — ВСЕГДА последняя строка
	err = http.ListenAndServe(":8000", rt)
	if err != nil {
		fmt.Println(err)
	}
}
