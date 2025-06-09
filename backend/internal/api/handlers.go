package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/stevebennett/slack-invite-mgr/backend/internal/config"
	"github.com/stevebennett/slack-invite-mgr/backend/internal/services"
)

// Invite represents a single invite from the spreadsheet
type Invite struct {
	Name            string `json:"name"`
	Role            string `json:"role"`
	Email           string `json:"email"`
	Company         string `json:"company"`
	YearsExperience string `json:"yearsExperience"`
	Reasons         string `json:"reasons"`
	Source          string `json:"source"`
}

// UpdateInviteStatusRequest represents the request to update invite statuses
type UpdateInviteStatusRequest struct {
	Emails []string `json:"emails"`
	Status string   `json:"status"`
}

// GetOutstandingInvitesHandler handles requests to get outstanding invites
func GetOutstandingInvitesHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Create sheets service
		sheetsService, err := services.NewSheetsService(r.Context(), &config.SheetsConfig{
			CredentialsFile: cfg.GoogleCredentialsFile,
			TokenFile:       cfg.GoogleTokenFile,
			SpreadsheetID:   cfg.GoogleSpreadsheetID,
			SheetName:       cfg.GoogleSheetName,
		})
		if err != nil {
			http.Error(w, "Failed to create sheets service", http.StatusInternalServerError)
			return
		}

		// Get sheet data
		data, err := sheetsService.GetSheetData(r.Context())
		if err != nil {
			http.Error(w, "Failed to get sheet data", http.StatusInternalServerError)
			return
		}

		// Convert data to invites
		var invites []Invite
		for _, row := range data {
			if len(row) < 9 {
				continue
			}

			invite := Invite{
				Name:            getString(row, 1), // Column B
				Role:            getString(row, 2), // Column C
				Email:           getString(row, 3), // Column D
				Company:         getString(row, 5), // Column F
				YearsExperience: getString(row, 6), // Column G
				Reasons:         getString(row, 7), // Column H
				Source:          getString(row, 8), // Column I
			}
			invites = append(invites, invite)
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Write response
		if err := json.NewEncoder(w).Encode(invites); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// UpdateInviteStatusHandler handles requests to update invite statuses
func UpdateInviteStatusHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse request body
		var req UpdateInviteStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Create sheets service
		sheetsService, err := services.NewSheetsService(r.Context(), &config.SheetsConfig{
			CredentialsFile: cfg.GoogleCredentialsFile,
			TokenFile:       cfg.GoogleTokenFile,
			SpreadsheetID:   cfg.GoogleSpreadsheetID,
			SheetName:       cfg.GoogleSheetName,
		})
		if err != nil {
			http.Error(w, "Failed to create sheets service", http.StatusInternalServerError)
			return
		}

		// Update the status for each email
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		if err := sheetsService.UpdateInviteStatus(r.Context(), req.Emails, req.Status, timestamp); err != nil {
			http.Error(w, "Failed to update invite statuses", http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

// Helper function to safely get string values from interface slice
func getString(row []interface{}, index int) string {
	if len(row) <= index {
		return ""
	}
	if str, ok := row[index].(string); ok {
		return str
	}
	return ""
}
