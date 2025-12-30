package todo_test

import (
	"context"
	"testing"

	"github.com/gemini/go-todo/internal/storage/memory"
	"github.com/gemini/go-todo/internal/todo"
)

func TestService_CreateTodo(t *testing.T) {
	repo := memory.NewRepo()
	service := todo.NewService(repo)
	ctx := context.Background()

	t.Run("creates a valid todo", func(t *testing.T) {
		title := "Test Todo"
		description := "Test Description"

		createdTodo, err := service.CreateTodo(ctx, title, description)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if createdTodo.ID == 0 {
			t.Errorf("expected ID to be set")
		}
		if createdTodo.Title != title {
			t.Errorf("expected title %q, got %q", title, createdTodo.Title)
		}
		if createdTodo.Completed {
			t.Errorf("expected completed to be false")
		}
	})

	t.Run("returns error for invalid todo", func(t *testing.T) {
		_, err := service.CreateTodo(ctx, "", "")
		if err != todo.ErrInvalid {
			t.Errorf("expected error %v, got %v", todo.ErrInvalid, err)
		}
	})
}

func TestService_UpdateTodo(t *testing.T) {
	repo := memory.NewRepo()
	service := todo.NewService(repo)
	ctx := context.Background()

	created, _ := service.CreateTodo(ctx, "Initial Title", "Initial Desc")

	t.Run("updates a todo", func(t *testing.T) {
		updated, err := service.UpdateTodo(ctx, created.ID, "Updated Title", "Updated Desc", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if updated.Title != "Updated Title" {
			t.Errorf("expected updated title, got %q", updated.Title)
		}
		if !updated.Completed {
			t.Errorf("expected completed to be true")
		}
	})

	t.Run("returns error for not found", func(t *testing.T) {
		_, err := service.UpdateTodo(ctx, 999, "title", "desc", false)
		if err != todo.ErrNotFound {
			t.Errorf("expected error %v, got %v", todo.ErrNotFound, err)
		}
	})
}
