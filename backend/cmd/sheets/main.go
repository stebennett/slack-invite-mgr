package main

import (
	"context"
	"fmt"
	"log"
	"time"

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
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if err := sheetsService.UpdateDuplicateRequests(ctx, timestamp); err != nil {
		// Send error email
		emailService := services.NewEmailService(sheetsCfg.EmailRecipient)
		if emailErr := emailService.SendEmail(ctx, "Error Updating Duplicate Requests", fmt.Sprintf("Error: %v", err)); emailErr != nil {
			log.Printf("Failed to send error email: %v", emailErr)
		}
		log.Fatalf("Failed to update duplicate requests: %v", err)
	}

	// Get new invites count
	newInvites, err := sheetsService.GetNewInvites(ctx)
	if err != nil {
		// Send error email
		emailService := services.NewEmailService(sheetsCfg.EmailRecipient)
		if emailErr := emailService.SendEmail(ctx, "Error Retrieving New Invites", fmt.Sprintf("Error: %v", err)); emailErr != nil {
			log.Printf("Failed to send error email: %v", emailErr)
		}
		log.Fatalf("Failed to get new invites: %v", err)
	}

	// Send success email if there are new invites
	if newInvites > 0 {
		emailService := services.NewEmailService(sheetsCfg.EmailRecipient)
		if err := emailService.SendEmail(ctx, "New Invites Need Processing", fmt.Sprintf("There are %d new invites that need processing.", newInvites)); err != nil {
			log.Printf("Failed to send success email: %v", err)
		}
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

	// Print summary information
	log.Printf("Number of new invites: %d", newInvites)
	log.Printf("Number of duplicates found: %d", duplicateCount)
}
