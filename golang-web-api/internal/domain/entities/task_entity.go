package entities

import "time"

// TaskEntity represents a task in the domain.
// It is a pure data structure with no knowledge of databases, HTTP, or any technology.
// The domain layer defines WHAT exists; other layers decide HOW to persist or transport it.
type TaskEntity struct {
	ID          int64
	Title       string
	Description string
	Status      string // "pending", "in_progress", "done"
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
