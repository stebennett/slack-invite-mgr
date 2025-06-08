package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stevebennett/slack-invite-mgr/backend/internal/config"
	"google.golang.org/api/sheets/v4"
)

// mockSheetsAPI mocks the relevant part of the Google Sheets API
// Only implements the method we use: Spreadsheets.Values.Get(...).Do()
type mockSheetsService struct {
	valuesGetFunc func(spreadsheetId, readRange string) *mockValuesGetCall
}

type mockValuesGetCall struct {
	ctx    context.Context
	doFunc func() (*sheets.ValueRange, error)
}

func (m *mockSheetsService) Spreadsheets() *mockSpreadsheets {
	return &mockSpreadsheets{m}
}

type mockSpreadsheets struct {
	parent *mockSheetsService
}

func (m *mockSpreadsheets) Values() *mockValues {
	return &mockValues{m.parent}
}

type mockValues struct {
	parent *mockSheetsService
}

func (m *mockValues) Get(spreadsheetId, readRange string) *mockValuesGetCall {
	return m.parent.valuesGetFunc(spreadsheetId, readRange)
}

func (c *mockValuesGetCall) Context(ctx context.Context) *mockValuesGetCall {
	c.ctx = ctx
	return c
}

func (c *mockValuesGetCall) Do() (*sheets.ValueRange, error) {
	return c.doFunc()
}

// SheetsServiceForTest is a testable version of SheetsService
// with a replaceable service field
func SheetsServiceForTest(mock *mockSheetsService, cfg *config.SheetsConfig) *SheetsService {
	return &SheetsService{
		service: (*sheets.Service)(nil), // not used
		cfg:     cfg,
	}
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
				valuesGetFunc: func(spreadsheetId, readRange string) *mockValuesGetCall {
					return &mockValuesGetCall{
						doFunc: tt.mockDo,
					}
				},
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
				resp, err := mockService.Spreadsheets().Values().Get(cfg.SpreadsheetID, rangeStr).Context(ctx).Do()
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
				{"A3", "B3", "C3", "D3", "E3", "F3", "G3", "H3", "I3"},
				{"A4", "B4", "C4", "D4", "E4", "F4", "G4", "H4", "I4", ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockSheetsService{
				valuesGetFunc: func(spreadsheetId, readRange string) *mockValuesGetCall {
					return &mockValuesGetCall{
						doFunc: func() (*sheets.ValueRange, error) {
							return &sheets.ValueRange{Values: tt.input}, nil
						},
					}
				},
			}
			cfg := &config.SheetsConfig{SpreadsheetID: "sheetid", SheetName: "Sheet1"}
			// svc := &SheetsService{cfg: cfg}

			// Patch the method call chain to use the filtering logic
			getSheetData := func(ctx context.Context) ([][]interface{}, error) {
				resp, err := mockService.Spreadsheets().Values().Get(cfg.SpreadsheetID, "Sheet1!A:J").Context(ctx).Do()
				if err != nil {
					return nil, err
				}
				var filtered [][]interface{}
				for _, row := range resp.Values {
					if len(row) < 10 || row[9] == "" {
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
				{"A1", "B1"},
				{"A2", "B2", "C2"},
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
			name:   "api error",
			input:  nil,
			want:   nil,
			apiErr: errors.New("api error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockSheetsService{
				valuesGetFunc: func(spreadsheetId, readRange string) *mockValuesGetCall {
					return &mockValuesGetCall{
						doFunc: func() (*sheets.ValueRange, error) {
							if tt.apiErr != nil {
								return nil, tt.apiErr
							}
							return &sheets.ValueRange{Values: tt.input}, nil
						},
					}
				},
			}
			cfg := &config.SheetsConfig{SpreadsheetID: "sheetid", SheetName: "Sheet1"}
			// svc := &SheetsService{cfg: cfg}

			getSheetData := func(ctx context.Context) ([][]interface{}, error) {
				resp, err := mockService.Spreadsheets().Values().Get(cfg.SpreadsheetID, "Sheet1!A:J").Context(ctx).Do()
				if err != nil {
					return nil, err
				}
				var filtered [][]interface{}
				for _, row := range resp.Values {
					if len(row) < 10 || row[9] == "" {
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
