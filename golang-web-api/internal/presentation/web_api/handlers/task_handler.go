package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appInterfaces "github.com/bling-lwsa/devpool-base-web-api/internal/app/interfaces"
	"github.com/bling-lwsa/devpool-base-web-api/internal/app/messages"
)

// TaskHandler handles HTTP requests related to tasks.
// It depends on the TaskServiceInterface (not the concrete service),
// following the Dependency Inversion principle.
type TaskHandler struct {
	service appInterfaces.TaskServiceInterface
}

// NewTaskHandler creates a TaskHandler with the given service.
func NewTaskHandler(service appInterfaces.TaskServiceInterface) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask handles POST /v1/tasks.
// It binds the JSON body to a CreateTaskRequest, delegates to the service,
// and returns the created task with HTTP 201.
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req messages.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// ListTasks handles GET /v1/tasks.
// It delegates to the service and returns the list of tasks with HTTP 200.
func (h *TaskHandler) ListTasks(c *gin.Context) {
	responses, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}
