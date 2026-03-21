package main

import (
	"os"
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
	appLogger.Info().Msg("Application initialized successfully")

	if err := a.Engine.Run(":8080"); err != nil {
		appLogger.Fatal().Err(err).Msg("Server failed to run")
	}
	appLogger.Info().Msg("Server stopped")
}
