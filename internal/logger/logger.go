package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// NewLogger initializes and returns a zerolog.Logger instance.
// It can be configured for console output (development) or JSON output (production).
func NewLogger(env string) zerolog.Logger {
	var logger zerolog.Logger

	if env == "development" {
		// Pretty console output for development
		output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		logger = zerolog.New(output).With().Timestamp().Logger()
	} else {
		// JSON output for production (e.g., for log aggregators)
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	// Set the logging level for this specific logger instance
	return logger.Level(zerolog.InfoLevel)
}
