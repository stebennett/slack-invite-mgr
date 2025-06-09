package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	SlackToken            string
	SlackChannelID        string
	DatabasePath          string
	GoogleCredentialsFile string
	GoogleTokenFile       string
	GoogleSpreadsheetID   string
	GoogleSheetName       string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	return &Config{
		SlackToken:            os.Getenv("SLACK_TOKEN"),
		SlackChannelID:        os.Getenv("SLACK_CHANNEL_ID"),
		DatabasePath:          os.Getenv("DATABASE_PATH"),
		GoogleCredentialsFile: os.Getenv("GOOGLE_CREDENTIALS_FILE"),
		GoogleTokenFile:       os.Getenv("GOOGLE_TOKEN_FILE"),
		GoogleSpreadsheetID:   os.Getenv("GOOGLE_SPREADSHEET_ID"),
		GoogleSheetName:       os.Getenv("GOOGLE_SHEET_NAME"),
	}, nil
}
