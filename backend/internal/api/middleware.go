package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const (
	// LoggerKey is the context key for the request-scoped logger
	LoggerKey contextKey = "logger"
	// RequestIDKey is the context key for the request ID
	RequestIDKey contextKey = "request_id"
)

// responseWriter wraps http.ResponseWriter to capture status code and size
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// LoggingMiddleware creates HTTP request/response logging middleware
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate or use existing request ID
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Create request-scoped logger
			reqLogger := logger.With(
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
			)

			// Add logger and request ID to context
			ctx := context.WithValue(r.Context(), LoggerKey, reqLogger)
			ctx = context.WithValue(ctx, RequestIDKey, requestID)

			// Wrap response writer
			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Add request ID to response headers
			w.Header().Set("X-Request-ID", requestID)

			// Log request start
			reqLogger.Info("request started")

			// Process request
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// Log request completion
			duration := time.Since(start)
			reqLogger.Info("request completed",
				slog.Int("status", wrapped.status),
				slog.Int("size", wrapped.size),
				slog.Duration("duration", duration),
			)
		})
	}
}

// LoggerFromContext retrieves the logger from context or returns the default
func LoggerFromContext(ctx context.Context, defaultLogger *slog.Logger) *slog.Logger {
	if logger, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
		return logger
	}
	return defaultLogger
}

// RequestIDFromContext retrieves the request ID from context
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
