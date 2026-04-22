package mocks

import (
	"context"

	"github.com/bling-lwsa/devpool-base-web-api/internal/domain/entities"
)

// TaskRepositoryMock is a manual mock that implements the domain TaskRepository interface.
//
// Instead of using a mock generation tool (like mockery), we create a struct with
// function fields. Each test can inject the exact behaviour it needs by assigning
// a closure. This approach is more explicit and great for learning how interfaces
// and testing work in Go.
//
// Example usage in a test:
//
//	repo := &mocks.TaskRepositoryMock{
//	    CreateFn: func(ctx context.Context, task *entities.TaskEntity) error {
//	        task.ID = 1
//	        return nil
//	    },
//	}
type TaskRepositoryMock struct {
	CreateFn func(ctx context.Context, task *entities.TaskEntity) error
	ListFn   func(ctx context.Context) ([]entities.TaskEntity, error)
}

func (m *TaskRepositoryMock) Create(ctx context.Context, task *entities.TaskEntity) error {
	return m.CreateFn(ctx, task)
}

func (m *TaskRepositoryMock) List(ctx context.Context) ([]entities.TaskEntity, error) {
	return m.ListFn(ctx)
}
