package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "user_api/cmd/docs"
	"user_api/internal/app"
	"user_api/internal/database"
	"user_api/internal/logger"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")

	// Determine environment for logging
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Initialize the application logger instance
	appLogger := logger.NewLogger(env)
	appLogger.Info().Msgf("Starting application in %s environment", env)

	database.Connect()
	if err := database.RunMigrations(); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret"
		appLogger.Warn().Msg("JWT_SECRET not set, using default secret.")
	}

	a, err := app.NewApp(jwtSecret, appLogger)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to create application")
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: a.Engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	appLogger.Info().Msgf("server started on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info().Msgf("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Info().Msgf("server shutdown failed: %v", err)
	}
	appLogger.Info().Msg("Server stopped")
}
