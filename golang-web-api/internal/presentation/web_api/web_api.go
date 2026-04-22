package web_api

import (
	"github.com/bling-lwsa/devpool-base-web-api/internal/presentation/web_api/routers"
)

// API is the top-level struct that aggregates everything needed to run the web API.
// In the corporate chassi this also holds a Logger, Tracer, and HealthChecker.
// Here we keep only the Router since that is the essential piece.
type API struct {
	Router *routers.APIRouter
}
