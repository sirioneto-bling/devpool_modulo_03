package routers

import "github.com/gin-gonic/gin"

func (r *APIRouter) registerHealthRoutes(group *gin.RouterGroup) {
	group.GET("/health", r.healthHandler.HealthCheck)
	group.GET("/livez", r.healthHandler.Livez)
	group.GET("/readyz", r.healthHandler.Readyz)
}
