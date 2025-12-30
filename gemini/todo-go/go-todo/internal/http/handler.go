package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gemini/go-todo/internal/todo"
	"github.com/go-chi/chi/v5"
)

// TodoService defines the interface for todo-related operations.
type TodoService interface {
	CreateTodo(ctx context.Context, title, description string) (*todo.Todo, error)
	ListTodos(ctx context.Context, completed *bool) ([]*todo.Todo, error)
	GetTodo(ctx context.Context, id int64) (*todo.Todo, error)
	UpdateTodo(ctx context.Context, id int64, title, description string, completed bool) (*todo.Todo, error)
	DeleteTodo(ctx context.Context, id int64) error
}

// Handler handles HTTP requests for todos.
type Handler struct {
	service TodoService
	logger  *slog.Logger
}

// NewHandler creates a new HTTP handler for todos.
func NewHandler(service TodoService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers the todo routes.
func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/api/todos", func(r chi.Router) {
		r.Post("/", h.createTodo)
		r.Get("/", h.listTodos)
		r.Get("/{id}", h.getTodo)
		r.Put("/{id}", h.updateTodo)
		r.Delete("/{id}", h.deleteTodo)
	})
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.JSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid_request"})
		return
	}

	createdTodo, err := h.service.CreateTodo(r.Context(), req.Title, req.Description)
	if errors.Is(err, todo.ErrInvalid) {
		h.JSON(w, r, http.StatusBadRequest, map[string]interface{}{"error": "validation_error", "details": map[string]string{"title": "is required"}})
		return
	}
	if err != nil {
		h.logger.Error("failed to create todo", "error", err)
		h.JSON(w, r, http.StatusInternalServerError, map[string]string{"error": "internal_error"})
		return
	}

	h.JSON(w, r, http.StatusCreated, createdTodo)
}

func (h *Handler) listTodos(w http.ResponseWriter, r *http.Request) {
	var completed *bool
	if completedStr := r.URL.Query().Get("completed"); completedStr != "" {
		c, err := strconv.ParseBool(completedStr)
		if err != nil {
			h.JSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid_query_param"})
			return
		}
		completed = &c
	}

	todos, err := h.service.ListTodos(r.Context(), completed)
	if err != nil {
		h.logger.Error("failed to list todos", "error", err)
		h.JSON(w, r, http.StatusInternalServerError, map[string]string{"error": "internal_error"})
		return
	}

	h.JSON(w, r, http.StatusOK, todos)
}

func (h *Handler) getTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid_id"})
		return
	}

	t, err := h.service.GetTodo(r.Context(), id)
	if errors.Is(err, todo.ErrNotFound) {
		h.JSON(w, r, http.StatusNotFound, map[string]string{"error": "not_found", "message": "todo not found"})
		return
	}
	if err != nil {
		h.logger.Error("failed to get todo", "error", err)
		h.JSON(w, r, http.StatusInternalServerError, map[string]string{"error": "internal_error"})
		return
	}

	h.JSON(w, r, http.StatusOK, t)
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid_id"})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.JSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid_request"})
		return
	}

	updatedTodo, err := h.service.UpdateTodo(r.Context(), id, req.Title, req.Description, req.Completed)
	if errors.Is(err, todo.ErrNotFound) {
		h.JSON(w, r, http.StatusNotFound, map[string]string{"error": "not_found", "message": "todo not found"})
		return
	}
	if errors.Is(err, todo.ErrInvalid) {
		h.JSON(w, r, http.StatusBadRequest, map[string]interface{}{"error": "validation_error", "details": map[string]string{"title": "is required"}})
		return
	}
	if err != nil {
		h.logger.Error("failed to update todo", "error", err)
		h.JSON(w, r, http.StatusInternalServerError, map[string]string{"error": "internal_error"})
		return
	}

	h.JSON(w, r, http.StatusOK, updatedTodo)
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid_id"})
		return
	}

	err = h.service.DeleteTodo(r.Context(), id)
	if errors.Is(err, todo.ErrNotFound) {
		h.JSON(w, r, http.StatusNotFound, map[string]string{"error": "not_found", "message": "todo not found"})
		return
	}
	if err != nil {
		h.logger.Error("failed to delete todo", "error", err)
		h.JSON(w, r, http.StatusInternalServerError, map[string]string{"error": "internal_error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// JSON writes a JSON response.
func (h *Handler) JSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			h.logger.Error("failed to write json response", "error", err)
		}
	}
}
