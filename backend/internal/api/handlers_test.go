package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stevebennett/slack-invite-mgr/backend/internal/config"
	"github.com/stevebennett/slack-invite-mgr/backend/internal/services"
)

// SheetsServiceInterface defines the methods we need from the sheets service
type SheetsServiceInterface interface {
	GetSheetData(ctx context.Context) ([][]interface{}, error)
	UpdateInviteStatus(ctx context.Context, emails []string, status string, timestamp string) error
	UpdateDuplicateRequests(ctx context.Context) error
	GetNewInvites(ctx context.Context) (int, error)
}

// mockSheetsService implements services.SheetsServiceInterface for testing
type mockSheetsService struct {
	updateStatusErr error
	data            [][]interface{}
}

func (m *mockSheetsService) GetSheetData(ctx context.Context) ([][]interface{}, error) {
	return m.data, nil
}

func (m *mockSheetsService) UpdateInviteStatus(ctx context.Context, emails []string, status string, timestamp string) error {
	return m.updateStatusErr
}

func (m *mockSheetsService) UpdateDuplicateRequests(ctx context.Context) error {
	return nil
}

func (m *mockSheetsService) GetNewInvites(ctx context.Context) (int, error) {
	return 0, nil
}

// NewSheetsService is a mock factory function
func NewSheetsService(ctx context.Context, cfg *config.SheetsConfig) (services.SheetsServiceInterface, error) {
	return &mockSheetsService{}, nil
}

func TestUpdateInviteStatusHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    UpdateInviteStatusRequest
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful update",
			requestBody: UpdateInviteStatusRequest{
				Emails: []string{"test@example.com"},
				Status: "sent",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty emails list",
			requestBody: UpdateInviteStatusRequest{
				Emails: []string{},
				Status: "sent",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "multiple emails",
			requestBody: UpdateInviteStatusRequest{
				Emails: []string{"test1@example.com", "test2@example.com"},
				Status: "denied",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "sheets service error",
			requestBody: UpdateInviteStatusRequest{
				Emails: []string{"test@example.com"},
				Status: "sent",
			},
			mockError:      errors.New("sheets service error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid request body",
			requestBody: UpdateInviteStatusRequest{
				Emails: []string{"invalid-email"},
				Status: "invalidstatus",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK, // We don't validate the status or email format in the handler
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test configuration
			cfg := &config.Config{
				GoogleCredentialsFile: "test-credentials.json",
				GoogleTokenFile:       "test-token.json",
				GoogleSpreadsheetID:   "test-spreadsheet-id",
				GoogleSheetName:       "test-sheet",
			}

			// Create request body
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPatch, "/api/invites", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create handler with mock service
			handler := UpdateInviteStatusHandler(cfg)

			// Serve request
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}

			// For successful requests, check response body
			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}
				if response["status"] != "success" {
					t.Errorf("handler returned unexpected response: got %v want %v",
						response["status"], "success")
				}
			}
		})
	}
}

func TestGetOutstandingInvitesHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockData       [][]interface{}
		mockError      error
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "successful fetch",
			mockData: [][]interface{}{
				{"", "John Doe", "Developer", "john@example.com", "", "Company", "5", "Reasons", "Source", "", ""},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "empty sheet",
			mockData:       [][]interface{}{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "sheets service error",
			mockData:       nil,
			mockError:      errors.New("sheets service error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
		{
			name: "incomplete row data",
			mockData: [][]interface{}{
				{"", "John Doe"}, // Incomplete row
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test configuration
			cfg := &config.Config{
				GoogleCredentialsFile: "test-credentials.json",
				GoogleTokenFile:       "test-token.json",
				GoogleSpreadsheetID:   "test-spreadsheet-id",
				GoogleSheetName:       "test-sheet",
			}

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/invites", nil)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create handler with mock service
			handler := GetOutstandingInvitesHandler(cfg)

			// Serve request
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}

			// For successful requests, check response body
			if tt.expectedStatus == http.StatusOK {
				var response []Invite
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}
				if len(response) != tt.expectedCount {
					t.Errorf("handler returned unexpected number of invites: got %v want %v",
						len(response), tt.expectedCount)
				}
			}
		})
	}
}
