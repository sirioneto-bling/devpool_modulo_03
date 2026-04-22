package services

import (
	"context"
	"errors"

	"github.com/bling-lwsa/devpool-base-web-api/internal/app/mappers"
	"github.com/bling-lwsa/devpool-base-web-api/internal/app/messages"
	domainInterfaces "github.com/bling-lwsa/devpool-base-web-api/internal/domain/interfaces"
)

// TaskService orchestrates task business logic.
// It depends on the domain's TaskRepository interface (not the MySQL implementation),
// which keeps the service decoupled from infrastructure details.
type TaskService struct {
	repo domainInterfaces.TaskRepository
}

// NewTaskService creates a TaskService with the given repository.
// This is the constructor pattern in Go -- since Go has no classes, we use
// a plain function that returns a pointer to the struct.
func NewTaskService(repo domainInterfaces.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(ctx context.Context, req messages.CreateTaskRequest) (*messages.TaskResponse, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	entity := mappers.ToEntity(req)

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	response := mappers.ToResponse(*entity)
	return &response, nil
}

func (s *TaskService) List(ctx context.Context) ([]messages.TaskResponse, error) {
	entityList, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return mappers.ToResponseList(entityList), nil
}
