package handlers_test

import (
	"context"

	"github.com/bling-lwsa/devpool-base-web-api/internal/app/messages"
)

type taskServiceMock struct {
	CreateFn func(ctx context.Context, req messages.CreateTaskRequest) (*messages.TaskResponse, error)
	ListFn   func(ctx context.Context) ([]messages.TaskResponse, error)
}

func (m *taskServiceMock) Create(ctx context.Context, req messages.CreateTaskRequest) (*messages.TaskResponse, error) {
	return m.CreateFn(ctx, req)
}

func (m *taskServiceMock) List(ctx context.Context) ([]messages.TaskResponse, error) {
	return m.ListFn(ctx)
}
