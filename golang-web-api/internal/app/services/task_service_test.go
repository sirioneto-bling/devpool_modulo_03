package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bling-lwsa/devpool-base-web-api/internal/app/messages"
	"github.com/bling-lwsa/devpool-base-web-api/internal/app/mocks"
	"github.com/bling-lwsa/devpool-base-web-api/internal/app/services"
	"github.com/bling-lwsa/devpool-base-web-api/internal/domain/entities"
)

// TestTaskService_Create uses table-driven tests -- the idiomatic Go pattern.
// Each test case is a struct in a slice; t.Run creates a named subtest per case.
// This makes it trivial to add new scenarios: just append another struct.
func TestTaskService_Create(t *testing.T) {
	tests := []struct {
		name    string
		request messages.CreateTaskRequest
		mockFn  func(ctx context.Context, task *entities.TaskEntity) error
		wantErr bool
	}{
		{
			name:    "success",
			request: messages.CreateTaskRequest{Title: "Study Go", Description: "Complete the tour of Go"},
			mockFn: func(_ context.Context, task *entities.TaskEntity) error {
				task.ID = 1
				task.CreatedAt = time.Now()
				task.UpdatedAt = time.Now()
				return nil
			},
			wantErr: false,
		},
		{
			name:    "empty title returns validation error",
			request: messages.CreateTaskRequest{Title: "", Description: "No title"},
			mockFn:  nil, // should not reach the repo
			wantErr: true,
		},
		{
			name:    "repository error is propagated",
			request: messages.CreateTaskRequest{Title: "Valid title"},
			mockFn: func(_ context.Context, _ *entities.TaskEntity) error {
				return errors.New("database connection lost")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.TaskRepositoryMock{CreateFn: tt.mockFn}
			svc := services.NewTaskService(repo)

			result, err := svc.Create(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "Study Go", result.Title)
				assert.Equal(t, "pending", result.Status)
			}
		})
	}
}

func TestTaskService_List(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		mockFn  func(ctx context.Context) ([]entities.TaskEntity, error)
		wantLen int
		wantErr bool
	}{
		{
			name: "success with results",
			mockFn: func(_ context.Context) ([]entities.TaskEntity, error) {
				return []entities.TaskEntity{
					{ID: 1, Title: "Task 1", Status: "pending", CreatedAt: now, UpdatedAt: now},
					{ID: 2, Title: "Task 2", Status: "done", CreatedAt: now, UpdatedAt: now},
				}, nil
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "empty list is not an error",
			mockFn: func(_ context.Context) ([]entities.TaskEntity, error) {
				return []entities.TaskEntity{}, nil
			},
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "repository error is propagated",
			mockFn: func(_ context.Context) ([]entities.TaskEntity, error) {
				return nil, errors.New("query timeout")
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.TaskRepositoryMock{ListFn: tt.mockFn}
			svc := services.NewTaskService(repo)

			result, err := svc.List(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}
