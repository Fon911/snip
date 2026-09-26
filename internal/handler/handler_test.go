package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"snap/example.com/internal/service"
	"snap/example.com/internal/storage"

	"github.com/go-chi/chi/v5"
)

func TestShorten(t *testing.T) {
	setupPool(t)
	rt := chi.NewRouter()
	rt.Post("/api/shorten", Shorten)
	t.Run("Сокрощение", func(t *testing.T) {
		body := strings.NewReader(`{"url": "https://www.youtube.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
		rec := httptest.NewRecorder()

		rt.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("статус: got %d, want 201; тело: %s", rec.Code, rec.Body.String())
		}

		if !strings.Contains(rec.Body.String(), `"code"`) {
			t.Errorf("в ответе нет id: %s", rec.Body.String())
		}

	})

	t.Run("Кривой JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader("musor"))
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("статус: got %d, want 400", rec.Code)
		}

	})

	t.Run("Пустой URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url": ""}`))
		rec := httptest.NewRecorder()
		rt.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("статус: got %d, want 400", rec.Code)
		}

	})

}

func setupPool(t *testing.T) {
	t.Helper()
	if Pool != nil {
		return
	}

	pool, err := storage.NewPool("postgres://postgres:secret@localhost:5400/postgres")
	if err != nil {
		t.Fatalf("база недоступна (docker start snipdb?): %v", err)
	}
	Pool = pool

}

func TestGet(t *testing.T) {
	setupPool(t)
	rt := chi.NewRouter()
	rt.Get("/{code}", Get)
	t.Run("Редирект", func(t *testing.T) {
		link, err := service.Shorten(Pool, "https://ya.ru")
		if err != nil {
			t.Fatalf("Ошибка Shorten: %v", err)
		}
		req := httptest.NewRequest(http.MethodGet, "/"+link, nil)
		rec := httptest.NewRecorder()

		rt.ServeHTTP(rec, req)
		if rec.Code != 302 {
			t.Errorf("статус: got %d, want 302", rec.Code)

		}

		loc := rec.Header().Get("Location")
		if loc != "https://ya.ru" {
			t.Errorf("Location: got %s, want %s", loc, "https://ya.ru")
		}

	})

	t.Run("Нет таког кода", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+"nope1234x", nil)
		rec := httptest.NewRecorder()

		rt.ServeHTTP(rec, req)
		if rec.Code != 404 {
			t.Errorf("статус: got %d, want 404", rec.Code)

		}

	})
}

func TestStats(t *testing.T) {
	setupPool(t)
	rt := chi.NewRouter()
	rt.Get("/api/stats/{code}", Stats)
	t.Run("Статистика", func(t *testing.T) {
		link, err := service.Shorten(Pool, "https://ya.ru")
		if err != nil {
			t.Fatalf("Ошибка Shorten: %v", err)
		}
		req := httptest.NewRequest(http.MethodGet, "/api/stats/"+link, nil)
		rec := httptest.NewRecorder()

		rt.ServeHTTP(rec, req)
		if !strings.Contains(rec.Body.String(), `"code"`) {
			t.Errorf("Ответ не полный: %s", rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), `"url"`) {
			t.Errorf("Ответ не полный: %s", rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), `"clicks"`) {
			t.Errorf("Ответ не полный: %s", rec.Body.String())
		}

	})

	t.Run("Нет таког кода", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/stats/"+"nope1234x", nil)
		rec := httptest.NewRecorder()

		rt.ServeHTTP(rec, req)
		if rec.Code != 404 {
			t.Errorf("статус: got %d, want 404", rec.Code)

		}

	})
}
