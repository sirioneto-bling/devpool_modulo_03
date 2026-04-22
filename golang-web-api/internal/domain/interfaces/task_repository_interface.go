package interfaces

import (
	"context"

	"github.com/bling-lwsa/devpool-base-web-api/internal/domain/entities"
)

// TaskRepository defines the contract for persisting tasks.
//
// This interface lives in the domain layer because the domain decides WHAT it needs,
// not HOW it is implemented. The infrastructure layer provides the concrete implementation
// (e.g. MySQL). This is Dependency Inversion -- the "D" in SOLID.
//
// In Go, interfaces are satisfied implicitly: any struct that has these methods
// automatically implements TaskRepository, without an explicit "implements" keyword.
//
// Exercise: add GetByID, Update and Delete methods to practice extending the contract.
type TaskRepository interface {
	Create(ctx context.Context, task *entities.TaskEntity) error
	List(ctx context.Context) ([]entities.TaskEntity, error)
}
