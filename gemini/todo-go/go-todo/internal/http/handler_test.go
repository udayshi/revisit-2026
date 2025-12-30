package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gemini/go-todo/internal/storage/memory"
	"github.com/gemini/go-todo/internal/todo"
	httpHandler "github.com/gemini/go-todo/internal/http"
	"github.com/go-chi/chi/v5"
)

func TestHandler_CreateTodo(t *testing.T) {
	repo := memory.NewRepo()
	service := todo.NewService(repo)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := httpHandler.NewHandler(service, logger)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	t.Run("creates a todo", func(t *testing.T) {
		body := `{"title": "Test", "description": "A test todo"}`
		req := httptest.NewRequest("POST", "/api/todos", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var createdTodo todo.Todo
		if err := json.NewDecoder(rr.Body).Decode(&createdTodo); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}

		if createdTodo.Title != "Test" {
			t.Errorf("expected title 'Test', got '%s'", createdTodo.Title)
		}
	})

	t.Run("returns bad request for invalid json", func(t *testing.T) {
		body := `{"title": "Test"`
		req := httptest.NewRequest("POST", "/api/todos", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestHandler_GetTodo(t *testing.T) {
	repo := memory.NewRepo()
	service := todo.NewService(repo)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := httpHandler.NewHandler(service, logger)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	// Pre-populate a todo
	created, _ := service.CreateTodo(context.Background(), "Existing Todo", "")

	t.Run("gets a todo", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/todos/"+strconv.FormatInt(created.ID, 10), nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("returns not found for non-existent todo", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/todos/999", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}
