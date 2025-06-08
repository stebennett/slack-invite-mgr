package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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

// getSheetIDByName fetches the SheetId for a given sheet name
func (s *SheetsService) getSheetIDByName(ctx context.Context, sheetName string) (int64, error) {
	spreadsheet, err := s.service.Spreadsheets.Get(s.cfg.SpreadsheetID).Context(ctx).Do()
	if err != nil {
		return 0, fmt.Errorf("failed to get spreadsheet metadata: %w", err)
	}
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties != nil && sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetId, nil
		}
	}
	return 0, fmt.Errorf("sheet with name '%s' not found", sheetName)
}

// UpdateDuplicateRequests checks for duplicate email addresses in column D and updates column J to "Duplicate" and column K with the current timestamp for duplicate rows.
func (s *SheetsService) UpdateDuplicateRequests(ctx context.Context) error {
	// Get the correct SheetId for the sheet name
	sheetId, err := s.getSheetIDByName(ctx, s.cfg.SheetName)
	if err != nil {
		return err
	}
	// Define the range to read (columns A-K)
	rangeStr := fmt.Sprintf("%s!A:K", s.cfg.SheetName)

	// Make the API call to get all rows
	resp, err := s.service.Spreadsheets.Values.Get(s.cfg.SpreadsheetID, rangeStr).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to retrieve sheet data: %w", err)
	}

	// Map to track the first occurrence of each email address
	emailMap := make(map[string]int)
	// Slice to hold the rows that need to be updated
	var updates []*sheets.ValueRange

	// Iterate through the rows
	for i, row := range resp.Values {
		if len(row) < 4 {
			continue // Skip rows with fewer than 4 columns
		}
		email, ok := row[3].(string)
		if !ok {
			continue // Skip if email is not a string
		}

		// Check if this email has been seen before
		if firstIndex, exists := emailMap[email]; exists {
			// If the first occurrence has an empty column J, mark this row as duplicate
			if len(resp.Values[firstIndex]) < 10 || resp.Values[firstIndex][9] == "" {
				// Ensure the row has at least 11 columns
				for len(row) < 11 {
					row = append(row, "")
				}
				// Update column J to "Duplicate" and column K with the current timestamp
				row[9] = "Duplicate"
				row[10] = time.Now().Format(time.RFC3339)
				// Prepare the update
				updates = append(updates, &sheets.ValueRange{
					Range:  fmt.Sprintf("%s!A%d:K%d", s.cfg.SheetName, i+1, i+1),
					Values: [][]interface{}{row},
				})
			} else {
				// If the first occurrence has a non-empty column J, mark this row as duplicate if its column J is empty
				if len(row) < 10 || row[9] == "" {
					for len(row) < 11 {
						row = append(row, "")
					}
					row[9] = "Duplicate"
					row[10] = time.Now().Format(time.RFC3339)
					updates = append(updates, &sheets.ValueRange{
						Range:  fmt.Sprintf("%s!A%d:K%d", s.cfg.SheetName, i+1, i+1),
						Values: [][]interface{}{row},
					})
				}
			}
		} else {
			// First occurrence of this email
			emailMap[email] = i
		}
	}

	// Apply the updates if any
	if len(updates) > 0 {
		// Prepare the batch update request
		var requests []*sheets.Request
		for _, update := range updates {
			// Parse the range string into a GridRange
			parts := strings.Split(update.Range, "!")
			if len(parts) != 2 {
				continue
			}
			cellRange := parts[1]
			// Assuming the range is in the format 'A1:K1', extract start and end
			startCell := strings.Split(cellRange, ":")[0]
			endCell := strings.Split(cellRange, ":")[1]
			startCol := startCell[0] - 'A'
			endCol := endCell[0] - 'A'
			startRow, _ := strconv.Atoi(startCell[1:])
			endRow, _ := strconv.Atoi(endCell[1:])

			// Convert update.Values to []*sheets.RowData
			var rows []*sheets.RowData
			for _, row := range update.Values {
				var cells []*sheets.CellData
				for _, cell := range row {
					cells = append(cells, &sheets.CellData{
						UserEnteredValue: &sheets.ExtendedValue{
							StringValue: func() *string { s := fmt.Sprintf("%v", cell); return &s }(),
						},
					})
				}
				rows = append(rows, &sheets.RowData{
					Values: cells,
				})
			}

			requests = append(requests, &sheets.Request{
				UpdateCells: &sheets.UpdateCellsRequest{
					Range: &sheets.GridRange{
						SheetId:          sheetId,
						StartRowIndex:    int64(startRow - 1),
						EndRowIndex:      int64(endRow),
						StartColumnIndex: int64(startCol),
						EndColumnIndex:   int64(endCol + 1),
					},
					Rows:   rows,
					Fields: "userEnteredValue",
				},
			})
		}

		// Execute the batch update
		_, err = s.service.Spreadsheets.BatchUpdate(s.cfg.SpreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}).Fields("*").Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("failed to update duplicate rows: %w", err)
		}
	}

	return nil
}
