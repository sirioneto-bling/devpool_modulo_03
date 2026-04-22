package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler provides health-check endpoints.
// It depends directly on *sql.DB (not on a service or repository) because
// health checks are an infrastructure concern, not a domain one.
// Not everything needs all DDD layers -- this is a good example of pragmatism.
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a HealthHandler with the given database connection.
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck pings the database and returns the overall system health.
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	err := h.db.PingContext(c.Request.Context())

	status := "healthy"
	httpStatus := http.StatusOK
	if err != nil {
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status": status,
		"checks": gin.H{
			"mysql": status,
		},
	})
}

// Livez returns 200 if the process is alive. Kubernetes uses this to know
// whether to restart the container.
func (h *HealthHandler) Livez(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

// Readyz checks if the application is ready to receive traffic.
// It verifies the database connection before reporting readiness.
func (h *HealthHandler) Readyz(c *gin.Context) {
	if err := h.db.PingContext(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
