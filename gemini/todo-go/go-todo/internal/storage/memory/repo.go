package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/gemini/go-todo/internal/todo"
)

// Repo is an in-memory implementation of the todo.Repository.
type Repo struct {
	mu    sync.RWMutex
	todos map[int64]*todo.Todo
	nextID int64
}

// NewRepo creates a new in-memory repository.
func NewRepo() *Repo {
	return &Repo{
		todos:  make(map[int64]*todo.Todo),
		nextID: 1,
	}
}

// Create creates a new todo.
func (r *Repo) Create(ctx context.Context, t *todo.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t.ID = r.nextID
	r.nextID++
	r.todos[t.ID] = t
	return nil
}

// FindAll returns all todos.
func (r *Repo) FindAll(ctx context.Context, completed *bool) ([]*todo.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*todo.Todo
	for _, t := range r.todos {
		if completed == nil || *completed == t.Completed {
			result = append(result, t)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

// FindByID finds a todo by its ID.
func (r *Repo) FindByID(ctx context.Context, id int64) (*todo.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.todos[id]
	if !ok {
		return nil, todo.ErrNotFound
	}
	return t, nil
}

// Update updates a todo.
func (r *Repo) Update(ctx context.Context, t *todo.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[t.ID]; !ok {
		return todo.ErrNotFound
	}
	r.todos[t.ID] = t
	return nil
}

// Delete deletes a todo by its ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[id]; !ok {
		return todo.ErrNotFound
	}
	delete(r.todos, id)
	return nil
}
