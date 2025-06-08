package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/stevebennett/slack-invite-mgr/backend/internal/config"
	"github.com/stevebennett/slack-invite-mgr/backend/internal/services"
)

func main() {
	// Load configuration
	sheetsCfg := config.LoadSheetsConfig()

	// Validate required configuration
	if sheetsCfg.CredentialsFile == "" {
		log.Fatal("GOOGLE_CREDENTIALS_FILE environment variable is required")
	}
	if sheetsCfg.SpreadsheetID == "" {
		log.Fatal("GOOGLE_SPREADSHEET_ID environment variable is required")
	}
	if sheetsCfg.SheetName == "" {
		log.Fatal("GOOGLE_SHEET_NAME environment variable is required")
	}

	// Create context
	ctx := context.Background()

	// Create sheets service
	sheetsService, err := services.NewSheetsService(ctx, sheetsCfg)
	if err != nil {
		log.Fatalf("Failed to create sheets service: %v", err)
	}

	// Get sheet data
	data, err := sheetsService.GetSheetData(ctx)
	if err != nil {
		log.Fatalf("Failed to get sheet data: %v", err)
	}

	// Print data in a readable format
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Fatalf("Failed to encode data: %v", err)
	}
}
