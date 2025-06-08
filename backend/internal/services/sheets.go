package services

import (
	"context"
	"fmt"

	"github.com/stevebennett/slack-invite-mgr/backend/internal/config"
	"google.golang.org/api/sheets/v4"
)

// SheetsService handles operations related to Google Sheets
type SheetsService struct {
	service *sheets.Service
	cfg     *config.SheetsConfig
}

// NewSheetsService creates a new SheetsService instance
func NewSheetsService(ctx context.Context, cfg *config.SheetsConfig) (*SheetsService, error) {
	service, err := config.GetSheetsService(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}

	return &SheetsService{
		service: service,
		cfg:     cfg,
	}, nil
}

// GetSheetData retrieves data from the specified sheet range
func (s *SheetsService) GetSheetData(ctx context.Context) ([][]interface{}, error) {
	// Define the range to read (columns A-J)
	rangeStr := fmt.Sprintf("%s!A:J", s.cfg.SheetName)

	// Make the API call
	resp, err := s.service.Spreadsheets.Values.Get(s.cfg.SpreadsheetID, rangeStr).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sheet data: %w", err)
	}

	// Filter rows where column J (index 9) is empty
	var filtered [][]interface{}
	for _, row := range resp.Values {
		if len(row) < 10 || row[9] == "" {
			filtered = append(filtered, row)
		}
	}

	return filtered, nil
}
