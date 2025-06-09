package api

import (
	"encoding/json"
	"net/http"

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
