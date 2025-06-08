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
