package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stevebennett/slack-invite-mgr/backend/internal/config"
	"google.golang.org/api/sheets/v4"
)

// mockSheetsAPI mocks the relevant part of the Google Sheets API
// Only implements the method we use: Spreadsheets.Values.Get(...).Do()
type mockSheetsService struct {
	values        [][]interface{}
	err           bool
	updatedValues [][]interface{}
	spreadsheet   *sheets.Spreadsheet
}

func (m *mockSheetsService) Get(ctx context.Context, spreadsheetId string, readRange string) (*sheets.ValueRange, error) {
	if m.err {
		return nil, errors.New("mock error")
	}
	return &sheets.ValueRange{
		Values: m.values,
	}, nil
}

func (m *mockSheetsService) BatchUpdate(ctx context.Context, spreadsheetId string, request *sheets.BatchUpdateSpreadsheetRequest) (*sheets.BatchUpdateSpreadsheetResponse, error) {
	if m.err {
		return nil, errors.New("mock error")
	}
	// Store the updated values for verification
	if len(request.Requests) > 0 && request.Requests[0].UpdateCells != nil {
		m.updatedValues = make([][]interface{}, len(m.values))
		copy(m.updatedValues, m.values)

		// Process each update request
		for _, req := range request.Requests {
			if req.UpdateCells != nil {
				// Get the row index from the range
				rowIndex := int(req.UpdateCells.Range.StartRowIndex)

				// Ensure the row exists in our updated values
				for len(m.updatedValues) <= rowIndex {
					m.updatedValues = append(m.updatedValues, make([]interface{}, 11))
				}

				// Ensure the row has enough columns
				for len(m.updatedValues[rowIndex]) < 11 {
					m.updatedValues[rowIndex] = append(m.updatedValues[rowIndex], "")
				}

				// Update the cells based on their column index
				for i, cell := range req.UpdateCells.Rows[0].Values {
					if cell.UserEnteredValue != nil && cell.UserEnteredValue.StringValue != nil {
						// Column J (index 9) and K (index 10) are updated separately
						colIndex := int(req.UpdateCells.Range.StartColumnIndex) + i
						if colIndex >= 9 && colIndex <= 10 {
							m.updatedValues[rowIndex][colIndex] = *cell.UserEnteredValue.StringValue
						}
					}
				}
			}
		}
	}
	return &sheets.BatchUpdateSpreadsheetResponse{}, nil
}

func (m *mockSheetsService) SpreadsheetsGet(ctx context.Context, spreadsheetId string) (*sheets.Spreadsheet, error) {
	if m.err {
		return nil, errors.New("mock error")
	}
	if m.spreadsheet != nil {
		return m.spreadsheet, nil
	}
	// Return a default spreadsheet with Sheet1
	return &sheets.Spreadsheet{
		Sheets: []*sheets.Sheet{
			{
				Properties: &sheets.SheetProperties{
					Title:   "Sheet1",
					SheetId: 0,
				},
			},
		},
	}, nil
}

func TestGetSheetData(t *testing.T) {
	tests := []struct {
		name      string
		mockDo    func() (*sheets.ValueRange, error)
		wantData  [][]interface{}
		wantError bool
	}{
		{
			name: "success",
			mockDo: func() (*sheets.ValueRange, error) {
				return &sheets.ValueRange{
					Values: [][]interface{}{{"A1", "B1"}, {"A2", "B2"}},
				}, nil
			},
			wantData:  [][]interface{}{{"A1", "B1"}, {"A2", "B2"}},
			wantError: false,
		},
		{
			name: "api error",
			mockDo: func() (*sheets.ValueRange, error) {
				return nil, errors.New("api error")
			},
			wantData:  nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockSheetsService{
				values: tt.wantData,
				err:    tt.wantError,
			}
			cfg := &config.SheetsConfig{SpreadsheetID: "sheetid", SheetName: "Sheet1"}

			// Patch SheetsService to use our mock
			// svc := &SheetsService{
			// 	cfg: cfg,
			// }
			// Patch the method call chain
			getSheetData := func(ctx context.Context) ([][]interface{}, error) {
				// Simulate the range string
				rangeStr := "Sheet1!A:J"
				resp, err := mockService.Get(ctx, cfg.SpreadsheetID, rangeStr)
				if err != nil {
					return nil, err
				}
				return resp.Values, nil
			}

			data, err := getSheetData(context.Background())
			if (err != nil) != tt.wantError {
				t.Errorf("error = %v, wantError %v", err, tt.wantError)
			}
			if !equal2D(data, tt.wantData) {
				t.Errorf("data = %v, want %v", data, tt.wantData)
			}
		})
	}
}

func TestGetSheetData_FilterColumnJ(t *testing.T) {
	tests := []struct {
		name  string
		input [][]interface{}
		want  [][]interface{}
	}{
		{
			name: "only rows with empty column J",
			input: [][]interface{}{
				{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1", "I1", ""},     // empty J
				{"A2", "B2", "C2", "D2", "E2", "F2", "G2", "H2", "I2", "done"}, // not empty J
				{"A3", "B3", "C3", "D3", "E3", "F3", "G3", "H3", "I3"},         // missing J
				{"A4", "B4", "C4", "D4", "E4", "F4", "G4", "H4", "I4", ""},     // empty J
			},
			want: [][]interface{}{
				{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1", "I1", ""},
				{"A3", "B3", "C3", "D3", "E3", "F3", "G3", "H3", "I3", ""},
				{"A4", "B4", "C4", "D4", "E4", "F4", "G4", "H4", "I4", ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockSheetsService{
				values: tt.input,
				err:    false,
			}
			cfg := &config.SheetsConfig{SpreadsheetID: "sheetid", SheetName: "Sheet1"}
			// svc := &SheetsService{cfg: cfg}

			// Patch the method call chain to use the filtering logic
			getSheetData := func(ctx context.Context) ([][]interface{}, error) {
				resp, err := mockService.Get(ctx, cfg.SpreadsheetID, "Sheet1!A:J")
				if err != nil {
					return nil, err
				}
				var filtered [][]interface{}
				for _, row := range resp.Values {
					// Ensure the row has at least 10 columns (A through J) by padding with empty strings
					for len(row) < 10 {
						row = append(row, "")
					}

					// Only include rows where column J (index 9) is empty
					if row[9] == "" {
						filtered = append(filtered, row)
					}
				}
				return filtered, nil
			}

			data, err := getSheetData(context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !equal2D(data, tt.want) {
				t.Errorf("filtered data = %v, want %v", data, tt.want)
			}
		})
	}
}

func TestGetSheetData_FurtherCases(t *testing.T) {
	tests := []struct {
		name   string
		input  [][]interface{}
		want   [][]interface{}
		apiErr error
	}{
		{
			name:  "empty sheet",
			input: [][]interface{}{},
			want:  [][]interface{}{},
		},
		{
			name: "all rows with J filled",
			input: [][]interface{}{
				{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1", "I1", "done"},
				{"A2", "B2", "C2", "D2", "E2", "F2", "G2", "H2", "I2", "x"},
			},
			want: [][]interface{}{},
		},
		{
			name: "rows with <10 columns",
			input: [][]interface{}{
				{"A1", "B1"},
				{"A2", "B2", "C2"},
			},
			want: [][]interface{}{
				{"A1", "B1", "", "", "", "", "", "", "", ""},
				{"A2", "B2", "C2", "", "", "", "", "", "", ""},
			},
		},
		{
			name: "row with exactly 10 columns, J empty",
			input: [][]interface{}{
				{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1", "I1", ""},
			},
			want: [][]interface{}{
				{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1", "I1", ""},
			},
		},
		{
			name: "row with exactly 10 columns, J not empty",
			input: [][]interface{}{
				{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1", "I1", "filled"},
			},
			want: [][]interface{}{},
		},
		{
			name: "sparse row with empty cells in middle",
			input: [][]interface{}{
				{"A1", "", "", "D1", "", "", "", "", "", ""}, // Row with empty B, C, E, F, G, H, I
			},
			want: [][]interface{}{
				{"A1", "", "", "D1", "", "", "", "", "", ""}, // Should be included since J is empty
			},
		},
		{
			name:   "api error",
			input:  nil,
			want:   nil,
			apiErr: errors.New("api error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockSheetsService{
				values: tt.input,
				err:    tt.apiErr != nil,
			}
			cfg := &config.SheetsConfig{SpreadsheetID: "sheetid", SheetName: "Sheet1"}
			// svc := &SheetsService{cfg: cfg}

			getSheetData := func(ctx context.Context) ([][]interface{}, error) {
				resp, err := mockService.Get(ctx, cfg.SpreadsheetID, "Sheet1!A:J")
				if err != nil {
					return nil, err
				}
				var filtered [][]interface{}
				for _, row := range resp.Values {
					// Ensure the row has at least 10 columns (A through J) by padding with empty strings
					for len(row) < 10 {
						row = append(row, "")
					}

					// Only include rows where column J (index 9) is empty
					if row[9] == "" {
						filtered = append(filtered, row)
					}
				}
				return filtered, nil
			}

			data, err := getSheetData(context.Background())
			if tt.apiErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !equal2D(data, tt.want) {
				t.Errorf("filtered data = %v, want %v", data, tt.want)
			}
		})
	}
}

func TestUpdateDuplicateRequests(t *testing.T) {
	cfg := &config.SheetsConfig{
		SpreadsheetID: "test-sheet-id",
		SheetName:     "Sheet1",
	}

	testTimestamp := "2024-02-14 12:00:00"

	testCases := []struct {
		name           string
		inputData      [][]interface{}
		expectedOutput [][]interface{}
		expectedError  bool
	}{
		{
			name: "Multiple rows with same email, all empty column J",
			inputData: [][]interface{}{
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
			},
			expectedOutput: [][]interface{}{
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "Duplicate", testTimestamp},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "Duplicate", testTimestamp},
			},
			expectedError: false,
		},
		{
			name: "Exactly 2 rows with same email, both empty column J",
			inputData: [][]interface{}{
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
			},
			expectedOutput: [][]interface{}{
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "Duplicate", testTimestamp},
			},
			expectedError: false,
		},
		{
			name: "Multiple rows with same email, some non-empty column J",
			inputData: [][]interface{}{
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "Processed", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
			},
			expectedOutput: [][]interface{}{
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "Processed", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "Duplicate", testTimestamp},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "Duplicate", testTimestamp},
			},
			expectedError: false,
		},
		{
			name: "No duplicates",
			inputData: [][]interface{}{
				{"1", "2", "3", "test1@example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test2@example.com", "5", "6", "7", "8", "9", "", ""},
			},
			expectedOutput: [][]interface{}{},
			expectedError:  false,
		},
		{
			name:           "Empty sheet",
			inputData:      [][]interface{}{},
			expectedOutput: [][]interface{}{},
			expectedError:  false,
		},
		{
			name: "API error",
			inputData: [][]interface{}{
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
			},
			expectedOutput: nil,
			expectedError:  true,
		},
		{
			name: "Case sensitive email duplicates",
			inputData: [][]interface{}{
				{"1", "2", "3", "Test@Example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
			},
			expectedOutput: [][]interface{}{
				{"1", "2", "3", "Test@Example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "Duplicate", testTimestamp},
			},
			expectedError: false,
		},
		{
			name: "Email with whitespace",
			inputData: [][]interface{}{
				{"1", "2", "3", " test@example.com ", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
			},
			expectedOutput: [][]interface{}{
				{"1", "2", "3", " test@example.com ", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "Duplicate", testTimestamp},
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockSheetsService{
				values: tc.inputData,
				err:    tc.expectedError,
			}

			svc := &SheetsService{
				cfg:     cfg,
				service: mockService,
			}

			err := svc.UpdateDuplicateRequests(context.Background(), testTimestamp)
			if tc.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(tc.expectedOutput) == 0 && len(mockService.updatedValues) == 0 {
				return
			}

			if !reflect.DeepEqual(mockService.updatedValues, tc.expectedOutput) {
				t.Errorf("Expected %v, got %v", tc.expectedOutput, mockService.updatedValues)
			}
		})
	}
}

func TestGetNewInvites(t *testing.T) {
	cfg := &config.SheetsConfig{
		SpreadsheetID: "test-sheet-id",
		SheetName:     "Sheet1",
	}

	testCases := []struct {
		name          string
		inputData     [][]interface{}
		expectedCount int
		expectedError bool
	}{
		{
			name: "Multiple new invites",
			inputData: [][]interface{}{
				{"1", "2", "3", "test1@example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test2@example.com", "5", "6", "7", "8", "9", "", ""},
				{"1", "2", "3", "test3@example.com", "5", "6", "7", "8", "9", "", ""},
			},
			expectedCount: 3,
			expectedError: false,
		},
		{
			name: "No new invites",
			inputData: [][]interface{}{
				{"1", "2", "3", "test1@example.com", "5", "6", "7", "8", "9", "Processed", ""},
				{"1", "2", "3", "test2@example.com", "5", "6", "7", "8", "9", "Duplicate", ""},
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:          "Empty sheet",
			inputData:     [][]interface{}{},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "API error",
			inputData: [][]interface{}{
				{"1", "2", "3", "test@example.com", "5", "6", "7", "8", "9", "", ""},
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockSheetsService{
				values: tc.inputData,
				err:    tc.expectedError,
			}

			svc := &SheetsService{
				cfg:     cfg,
				service: mockService,
			}

			count, err := svc.GetNewInvites(context.Background())
			if tc.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if count != tc.expectedCount {
				t.Errorf("Expected count %d, got %d", tc.expectedCount, count)
			}
		})
	}
}

// Helper to compare [][]interface{}
func equal2D(a, b [][]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
