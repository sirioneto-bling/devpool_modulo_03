package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bling-lwsa/devpool-base-web-api/internal/domain/entities"
)

// TaskRepositoryMySQL implements the domain TaskRepository interface using database/sql.
// It translates domain operations into SQL queries against a MySQL database.
type TaskRepositoryMySQL struct {
	db *sql.DB
}

// NewTaskRepositoryMySQL creates a new repository backed by the given database connection.
func NewTaskRepositoryMySQL(db *sql.DB) *TaskRepositoryMySQL {
	return &TaskRepositoryMySQL{db: db}
}

func (r *TaskRepositoryMySQL) Create(ctx context.Context, task *entities.TaskEntity) error {
	query := "INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)"

	result, err := r.db.ExecContext(ctx, query, task.Title, task.Description, task.Status)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	task.ID = id
	return nil
}

func (r *TaskRepositoryMySQL) List(ctx context.Context) ([]entities.TaskEntity, error) {
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []entities.TaskEntity
	for rows.Next() {
		var t entities.TaskEntity
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task rows: %w", err)
	}

	return tasks, nil
}
