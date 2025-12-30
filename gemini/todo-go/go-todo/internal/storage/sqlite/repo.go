package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gemini/go-todo/internal/todo"
	_ "github.com/mattn/go-sqlite3"
)

// Repo is a SQLite implementation of the todo.Repository.
type Repo struct {
	db *sql.DB
}

// NewRepo creates a new SQLite repository. It also runs migrations.
func NewRepo(dsn string) (*Repo, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Repo{db: db}, nil
}

// Close closes the database connection.
func (r *Repo) Close() error {
	return r.db.Close()
}

func runMigrations(db *sql.DB) error {
	migration, err := os.ReadFile("migrations/001_create_todos.sql")
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(migration)); err != nil {
		return err
	}
	return nil
}

// Create creates a new todo.
func (r *Repo) Create(ctx context.Context, t *todo.Todo) error {
	query := `INSERT INTO todos (title, description, completed, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, t.Title, t.Description, t.Completed, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = id
	return nil
}

// FindAll returns all todos.
func (r *Repo) FindAll(ctx context.Context, completed *bool) ([]*todo.Todo, error) {
	query := "SELECT id, title, description, completed, created_at, updated_at FROM todos"
	var args []interface{}
	if completed != nil {
		query += " WHERE completed = ?"
		args = append(args, *completed)
	}
	query += " ORDER BY id"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*todo.Todo
	for rows.Next() {
		t := &todo.Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

// FindByID finds a todo by its ID.
func (r *Repo) FindByID(ctx context.Context, id int64) (*todo.Todo, error) {
	query := "SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	t := &todo.Todo{}
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, todo.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Update updates a todo.
func (r *Repo) Update(ctx context.Context, t *todo.Todo) error {
	query := "UPDATE todos SET title = ?, description = ?, completed = ?, updated_at = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, t.Title, t.Description, t.Completed, t.UpdatedAt, t.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a todo by its ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM todos WHERE id = ?"
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return todo.ErrNotFound
	}
	return nil
}
