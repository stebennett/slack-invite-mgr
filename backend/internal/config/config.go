package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	GoogleCredentialsFile string
	GoogleTokenFile       string
	GoogleSpreadsheetID   string
	GoogleSheetName       string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	return &Config{
		GoogleCredentialsFile: os.Getenv("GOOGLE_CREDENTIALS_FILE"),
		GoogleTokenFile:       os.Getenv("GOOGLE_TOKEN_FILE"),
		GoogleSpreadsheetID:   os.Getenv("GOOGLE_SPREADSHEET_ID"),
		GoogleSheetName:       os.Getenv("GOOGLE_SHEET_NAME"),
	}, nil
}
