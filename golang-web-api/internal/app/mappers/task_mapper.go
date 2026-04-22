package mappers

import (
	"github.com/bling-lwsa/devpool-base-web-api/internal/app/messages"
	"github.com/bling-lwsa/devpool-base-web-api/internal/domain/entities"
)

// ToEntity converts a CreateTaskRequest into a domain TaskEntity.
// The mapper centralises the conversion so handlers and services stay clean.
func ToEntity(req messages.CreateTaskRequest) *entities.TaskEntity {
	return &entities.TaskEntity{
		Title:       req.Title,
		Description: req.Description,
		Status:      "pending",
	}
}

// ToResponse converts a domain TaskEntity into an API-friendly TaskResponse.
func ToResponse(entity entities.TaskEntity) messages.TaskResponse {
	return messages.TaskResponse{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		Status:      entity.Status,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// ToResponseList converts a slice of entities into a slice of responses.
func ToResponseList(entityList []entities.TaskEntity) []messages.TaskResponse {
	responses := make([]messages.TaskResponse, 0, len(entityList))
	for _, e := range entityList {
		responses = append(responses, ToResponse(e))
	}
	return responses
}
