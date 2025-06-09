package api

import (
	"net/http"

	"github.com/stevebennett/slack-invite-mgr/backend/internal/config"
)

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Invites endpoints
	mux.HandleFunc("/api/invites", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetOutstandingInvitesHandler(cfg)(w, r)
		case http.MethodPatch:
			UpdateInviteStatusHandler(cfg)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
