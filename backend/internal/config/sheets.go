package config

import (
	"context"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// SheetsConfig holds Google Sheets specific configuration
type SheetsConfig struct {
	CredentialsFile string
	TokenFile       string
	SpreadsheetID   string
	SheetName       string
	EmailRecipient  string
	EmailTemplate   string
}

// LoadSheetsConfig loads Google Sheets configuration from environment variables
func LoadSheetsConfig() *SheetsConfig {
	return &SheetsConfig{
		CredentialsFile: os.Getenv("GOOGLE_CREDENTIALS_FILE"),
		TokenFile:       os.Getenv("GOOGLE_TOKEN_FILE"),
		SpreadsheetID:   os.Getenv("GOOGLE_SPREADSHEET_ID"),
		SheetName:       os.Getenv("GOOGLE_SHEET_NAME"),
		EmailRecipient:  os.Getenv("EMAIL_RECIPIENT"),
		EmailTemplate:   os.Getenv("EMAIL_TEMPLATE_PATH"),
	}
}

// GetSheetsService creates a new Google Sheets service client
func GetSheetsService(ctx context.Context, cfg *SheetsConfig) (*sheets.Service, error) {
	// Read credentials file
	credentials, err := os.ReadFile(cfg.CredentialsFile)
	if err != nil {
		return nil, err
	}

	// Parse credentials
	config, err := google.JWTConfigFromJSON(credentials, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	// Create client
	client := config.Client(ctx)

	// Create service
	service, err := sheets.New(client)
	if err != nil {
		return nil, err
	}

	return service, nil
}
