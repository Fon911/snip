package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"snap/example.com/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

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

	code, err := service.Shorten(Pool, req.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "не смог сохранить")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ShortenResponse{Code: code})
}

func Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	url, err :=
		service.Resolve(Pool, code)
	if errors.Is(err, service.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "нет такого кода")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Ошибка сервера")
		return
	}

	http.Redirect(w, r, url, http.StatusFound)

}

func Stats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	url, clicks, err := service.Stats(Pool, code)
	if errors.Is(err, service.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "нет такой ссылки")
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
