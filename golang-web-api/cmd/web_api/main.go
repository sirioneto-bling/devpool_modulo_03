package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/bling-lwsa/devpool-base-web-api/internal/app/services"
	"github.com/bling-lwsa/devpool-base-web-api/internal/infrastructure/config"
	"github.com/bling-lwsa/devpool-base-web-api/internal/infrastructure/mysql"
	"github.com/bling-lwsa/devpool-base-web-api/internal/infrastructure/mysql/repositories"
	"github.com/bling-lwsa/devpool-base-web-api/internal/presentation/web_api/handlers"
	"github.com/bling-lwsa/devpool-base-web-api/internal/presentation/web_api/routers"
)

func main() {
	// ---- 1. Load .env (ignored in production where env vars come from the orchestrator) ----
	godotenv.Load()

	// ---- 2. Configuration ----
	cfg := config.LoadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// ---- 3. Infrastructure: database connection ----
	db, err := mysql.NewConnection(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("database connection established")

	// ---- 4. Manual Dependency Injection ----
	//
	// In the corporate chassi we use Google Wire to generate this wiring automatically.
	// Here we do it by hand so you can see exactly what depends on what:
	//
	//   main -> repository (infra) -> service (app) -> handler (presentation) -> router
	//
	// Each layer receives only the interfaces it needs. The concrete types are decided here.

	// Infrastructure
	taskRepo := repositories.NewTaskRepositoryMySQL(db)

	// Application
	taskService := services.NewTaskService(taskRepo)

	// Presentation
	healthHandler := handlers.NewHealthHandler(db)
	taskHandler := handlers.NewTaskHandler(taskService)

	engine := gin.Default()
	router := routers.NewRouter(engine, healthHandler, taskHandler)
	ginEngine := router.RegisterRoutes()

	// ---- 5. HTTP server with graceful shutdown ----
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: ginEngine,
	}

	// Start server in a goroutine so it doesn't block the main goroutine.
	// The main goroutine will wait for an interrupt signal below.
	go func() {
		slog.Info("starting server", "port", cfg.Port, "app", cfg.AppName)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal (Ctrl+C or SIGTERM from Kubernetes/ECS).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh

	slog.Warn("shutting down gracefully...")

	// Give in-flight requests up to 10 seconds to finish.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}

	slog.Info("server stopped")
}
