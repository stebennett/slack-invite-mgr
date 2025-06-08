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

	// Update duplicate requests
	if err := sheetsService.UpdateDuplicateRequests(ctx); err != nil {
		log.Fatalf("Failed to update duplicate requests: %v", err)
	}

	// Get updated sheet data to count duplicates
	updatedData, err := sheetsService.GetSheetData(ctx)
	if err != nil {
		log.Fatalf("Failed to get updated sheet data: %v", err)
	}

	// Count duplicates
	duplicateCount := 0
	for _, row := range updatedData {
		if len(row) >= 10 && row[9] == "Duplicate" {
			duplicateCount++
		}
	}

	// Print data in a readable format
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(updatedData); err != nil {
		log.Fatalf("Failed to encode data: %v", err)
	}

	// Print number of duplicates found
	log.Printf("Number of duplicates found: %d", duplicateCount)
}
