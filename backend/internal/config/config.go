package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	SlackToken     string
	SlackChannelID string
	DatabasePath   string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	return &Config{
		SlackToken:     os.Getenv("SLACK_TOKEN"),
		SlackChannelID: os.Getenv("SLACK_CHANNEL_ID"),
		DatabasePath:   os.Getenv("DATABASE_PATH"),
	}, nil
}
