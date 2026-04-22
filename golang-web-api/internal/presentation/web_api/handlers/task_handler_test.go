package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/bling-lwsa/devpool-base-web-api/internal/app/messages"
	"github.com/bling-lwsa/devpool-base-web-api/internal/presentation/web_api/handlers"
)

func setupTaskRouter(svc *taskServiceMock) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	handler := handlers.NewTaskHandler(svc)

	v1 := engine.Group("/v1")
	v1.POST("/tasks", handler.CreateTask)
	v1.GET("/tasks", handler.ListTasks)

	return engine
}

func TestCreateTask_Success(t *testing.T) {
	now := time.Now()
	svc := &taskServiceMock{
		CreateFn: func(_ context.Context, req messages.CreateTaskRequest) (*messages.TaskResponse, error) {
			return &messages.TaskResponse{
				ID:          1,
				Title:       req.Title,
				Description: req.Description,
				Status:      "pending",
				CreatedAt:   now,
				UpdatedAt:   now,
			}, nil
		},
	}

	router := setupTaskRouter(svc)

	body := `{"title":"Study Go","description":"Complete the tour of Go"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp messages.TaskResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.ID)
	assert.Equal(t, "Study Go", resp.Title)
	assert.Equal(t, "Complete the tour of Go", resp.Description)
	assert.Equal(t, "pending", resp.Status)
}

func TestCreateTask_InvalidJSON(t *testing.T) {
	svc := &taskServiceMock{}
	router := setupTaskRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "invalid")
}

func TestCreateTask_MissingTitle(t *testing.T) {
	svc := &taskServiceMock{}
	router := setupTaskRouter(svc)

	body := `{"description":"No title provided"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateTask_ServiceError(t *testing.T) {
	svc := &taskServiceMock{
		CreateFn: func(_ context.Context, _ messages.CreateTaskRequest) (*messages.TaskResponse, error) {
			return nil, errors.New("database connection lost")
		},
	}
	router := setupTaskRouter(svc)

	body := `{"title":"Valid title"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "database connection lost", resp["error"])
}

func TestListTasks_Success(t *testing.T) {
	now := time.Now()
	svc := &taskServiceMock{
		ListFn: func(_ context.Context) ([]messages.TaskResponse, error) {
			return []messages.TaskResponse{
				{ID: 1, Title: "Task 1", Status: "pending", CreatedAt: now, UpdatedAt: now},
				{ID: 2, Title: "Task 2", Status: "done", CreatedAt: now, UpdatedAt: now},
			}, nil
		},
	}
	router := setupTaskRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []messages.TaskResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "Task 1", resp[0].Title)
	assert.Equal(t, "Task 2", resp[1].Title)
}

func TestListTasks_Empty(t *testing.T) {
	svc := &taskServiceMock{
		ListFn: func(_ context.Context) ([]messages.TaskResponse, error) {
			return []messages.TaskResponse{}, nil
		},
	}
	router := setupTaskRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []messages.TaskResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Empty(t, resp)
}

func TestListTasks_ServiceError(t *testing.T) {
	svc := &taskServiceMock{
		ListFn: func(_ context.Context) ([]messages.TaskResponse, error) {
			return nil, errors.New("query timeout")
		},
	}
	router := setupTaskRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "query timeout", resp["error"])
}
