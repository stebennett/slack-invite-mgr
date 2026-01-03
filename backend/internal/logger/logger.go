package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Config holds logger configuration
type Config struct {
	Level   slog.Level
	AppName string
	Output  io.Writer
}

// ParseLevel converts a string level to slog.Level
func ParseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// New creates a new configured logger
func New(cfg Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	handler := slog.NewJSONHandler(cfg.Output, opts)

	// Add static fields for all log entries
	return slog.New(handler).With(
		slog.String("app", cfg.AppName),
	)
}

// FromEnv creates a logger configured from environment variables
func FromEnv(appName string) *slog.Logger {
	level := ParseLevel(os.Getenv("LOG_LEVEL"))
	return New(Config{
		Level:   level,
		AppName: appName,
		Output:  os.Stdout,
	})
}

// WithRequestID adds request ID to logger
func WithRequestID(logger *slog.Logger, requestID string) *slog.Logger {
	return logger.With(slog.String("request_id", requestID))
}
