package todo

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNotFound is returned when a todo is not found.
	ErrNotFound = errors.New("todo not found")
	// ErrInvalid is returned when a todo is invalid.
	ErrInvalid = errors.New("todo invalid")
)

// Repository defines the interface for todo storage.
type Repository interface {
	Create(ctx context.Context, todo *Todo) error
	FindAll(ctx context.Context, completed *bool) ([]*Todo, error)
	FindByID(ctx context.Context, id int64) (*Todo, error)
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id int64) error
}

// Service provides todo-related operations.
type Service struct {
	repo Repository
}

// NewService creates a new todo service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateTodo creates a new todo.
func (s *Service) CreateTodo(ctx context.Context, title, description string) (*Todo, error) {
	now := time.Now()
	todo := &Todo{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := todo.Validate(); err != nil {
		return nil, ErrInvalid
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// ListTodos lists all todos.
func (s *Service) ListTodos(ctx context.Context, completed *bool) ([]*Todo, error) {
	return s.repo.FindAll(ctx, completed)
}

// GetTodo gets a single todo by its ID.
func (s *Service) GetTodo(ctx context.Context, id int64) (*Todo, error) {
	return s.repo.FindByID(ctx, id)
}

// UpdateTodo updates a todo.
func (s *Service) UpdateTodo(ctx context.Context, id int64, title, description string, completed bool) (*Todo, error) {
	todo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	todo.Title = title
	todo.Description = description
	todo.Completed = completed
	todo.UpdatedAt = time.Now()

	if err := todo.Validate(); err != nil {
		return nil, ErrInvalid
	}

	if err := s.repo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// DeleteTodo deletes a todo by its ID.
func (s *Service) DeleteTodo(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
