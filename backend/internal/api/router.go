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

	// TODO: Add more routes as needed

	return mux
}
