package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/stevebennett/slack-invite-mgr/backend/internal/config"
	"github.com/stevebennett/slack-invite-mgr/backend/internal/logger"
	"github.com/stevebennett/slack-invite-mgr/backend/internal/services"
)

func main() {
	// Initialize logger
	log := logger.FromEnv("slack-invite-sheets")

	// Load configuration
	sheetsCfg := config.LoadSheetsConfig()

	// Validate required configuration
	if sheetsCfg.CredentialsFile == "" {
		log.Error("missing required configuration", slog.String("field", "GOOGLE_CREDENTIALS_FILE"))
		os.Exit(1)
	}
	if sheetsCfg.SpreadsheetID == "" {
		log.Error("missing required configuration", slog.String("field", "GOOGLE_SPREADSHEET_ID"))
		os.Exit(1)
	}
	if sheetsCfg.SheetName == "" {
		log.Error("missing required configuration", slog.String("field", "GOOGLE_SHEET_NAME"))
		os.Exit(1)
	}

	// Create context
	ctx := context.Background()

	// Create sheets service
	log.Info("creating sheets service")
	sheetsService, err := services.NewSheetsService(ctx, sheetsCfg)
	if err != nil {
		log.Error("failed to create sheets service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Update duplicate requests
	log.Info("updating duplicate requests")
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if err := sheetsService.UpdateDuplicateRequests(ctx, timestamp); err != nil {
		log.Error("failed to update duplicate requests", slog.String("error", err.Error()))
		// Send error email
		emailService := services.NewEmailService(sheetsCfg.EmailRecipient, sheetsCfg.EmailTemplate)
		if emailErr := emailService.SendEmail(ctx, "Error Updating Duplicate Requests", fmt.Sprintf("Error: %v", err)); emailErr != nil {
			log.Error("failed to send error email", slog.String("error", emailErr.Error()))
		}
		os.Exit(1)
	}

	// Get new invites count
	log.Info("retrieving new invites count")
	newInvites, err := sheetsService.GetNewInvites(ctx)
	if err != nil {
		log.Error("failed to get new invites", slog.String("error", err.Error()))
		// Send error email
		emailService := services.NewEmailService(sheetsCfg.EmailRecipient, sheetsCfg.EmailTemplate)
		if emailErr := emailService.SendEmail(ctx, "Error Retrieving New Invites", fmt.Sprintf("Error: %v", err)); emailErr != nil {
			log.Error("failed to send error email", slog.String("error", emailErr.Error()))
		}
		os.Exit(1)
	}

	// Send success email if there are new invites
	if newInvites > 0 {
		log.Info("sending notification email", slog.Int("new_invites", newInvites))
		emailService := services.NewEmailService(sheetsCfg.EmailRecipient, sheetsCfg.EmailTemplate)
		if err := emailService.SendEmail(ctx, "New Invites Need Processing", fmt.Sprintf("There are %d new invites that need processing.", newInvites)); err != nil {
			log.Error("failed to send success email", slog.String("error", err.Error()))
		}
	}

	// Get updated sheet data to count duplicates
	log.Debug("retrieving updated sheet data")
	updatedData, err := sheetsService.GetSheetData(ctx)
	if err != nil {
		log.Error("failed to get updated sheet data", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Count duplicates
	duplicateCount := 0
	for _, row := range updatedData {
		if len(row) >= 10 && row[9] == "Duplicate" {
			duplicateCount++
		}
	}

	// Log summary information
	log.Info("sheets sync completed",
		slog.Int("new_invites", newInvites),
		slog.Int("duplicates_found", duplicateCount),
	)
}
