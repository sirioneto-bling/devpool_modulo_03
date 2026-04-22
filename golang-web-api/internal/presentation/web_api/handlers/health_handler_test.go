package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/bling-lwsa/devpool-base-web-api/internal/presentation/web_api/handlers"
)

func setupHealthRouter(handler *handlers.HealthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	v1 := engine.Group("/v1")
	v1.GET("/health", handler.HealthCheck)
	v1.GET("/livez", handler.Livez)
	v1.GET("/readyz", handler.Readyz)

	return engine
}

func TestHealthCheck_Healthy(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing()

	handler := handlers.NewHealthHandler(db)
	router := setupHealthRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", resp["status"])

	checks := resp["checks"].(map[string]interface{})
	assert.Equal(t, "healthy", checks["mysql"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealthCheck_Unhealthy(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing().WillReturnError(assert.AnError)

	handler := handlers.NewHealthHandler(db)
	router := setupHealthRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "unhealthy", resp["status"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLivez(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handlers.NewHealthHandler(db)
	router := setupHealthRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/v1/livez", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "alive", resp["status"])
}

func TestReadyz_Ready(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing()

	handler := handlers.NewHealthHandler(db)
	router := setupHealthRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/v1/readyz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ready", resp["status"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReadyz_NotReady(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing().WillReturnError(assert.AnError)

	handler := handlers.NewHealthHandler(db)
	router := setupHealthRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/v1/readyz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)

	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "not ready", resp["status"])
	assert.NoError(t, mock.ExpectationsWereMet())
}
