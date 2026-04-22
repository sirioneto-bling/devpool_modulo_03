package interfaces

import (
	"context"

	"github.com/bling-lwsa/devpool-base-web-api/internal/app/messages"
)

// TaskServiceInterface defines the contract for the task business logic.
// The handler depends on this interface (not on the concrete service),
// which makes it easy to test or swap implementations.
type TaskServiceInterface interface {
	Create(ctx context.Context, req messages.CreateTaskRequest) (*messages.TaskResponse, error)
	List(ctx context.Context) ([]messages.TaskResponse, error)
}
