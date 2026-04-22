package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/bling-lwsa/devpool-base-web-api/internal/presentation/web_api/handlers"
)

// APIRouter groups all route registration in one place.
// In the corporate chassi this struct carries dozens of middlewares (auth, tracing,
// i18n, error handling). Here we keep it minimal so the routing logic is clear.
type APIRouter struct {
	engine        *gin.Engine
	healthHandler *handlers.HealthHandler
	taskHandler   *handlers.TaskHandler
}

// NewRouter creates a new APIRouter.
func NewRouter(
	engine *gin.Engine,
	healthHandler *handlers.HealthHandler,
	taskHandler *handlers.TaskHandler,
) *APIRouter {
	return &APIRouter{
		engine:        engine,
		healthHandler: healthHandler,
		taskHandler:   taskHandler,
	}
}

// RegisterRoutes wires every route group and returns the configured engine.
func (r *APIRouter) RegisterRoutes() *gin.Engine {
	r.engine.GET("/swagger/*any", func(c *gin.Context) {
		if c.Param("any") == "/" {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
			return
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	v1 := r.engine.Group("/v1")

	r.registerHealthRoutes(v1)
	r.registerTaskRoutes(v1)

	return r.engine
}
