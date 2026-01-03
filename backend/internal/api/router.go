package api

import (
	"log/slog"
	"net/http"

	"github.com/stevebennett/slack-invite-mgr/backend/internal/config"
)

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(cfg *config.Config, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint (no logging to reduce noise)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Invites endpoints
	mux.HandleFunc("/api/invites", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetOutstandingInvitesHandler(cfg, logger)(w, r)
		case http.MethodPatch:
			UpdateInviteStatusHandler(cfg, logger)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Frontend logs endpoint
	mux.HandleFunc("/api/logs", FrontendLogsHandler(logger))

	// Apply logging middleware
	return LoggingMiddleware(logger)(mux)
}
