package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/stevebennett/slack-invite-mgr/backend/internal/api"
	"github.com/stevebennett/slack-invite-mgr/backend/internal/config"
	"github.com/stevebennett/slack-invite-mgr/backend/internal/logger"
)

func main() {
	// Initialize logger
	log := logger.FromEnv("slack-invite-api")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize router
	router := api.NewRouter(cfg, log)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Info("server starting", slog.String("port", port))
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Error("server failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
