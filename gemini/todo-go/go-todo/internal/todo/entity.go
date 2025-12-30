package todo

import (
	"errors"
	"time"
)

// Todo represents a single todo item.
type Todo struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate validates the Todo struct.
func (t *Todo) Validate() error {
	if t.Title == "" {
		return errors.New("title is required")
	}
	return nil
}
