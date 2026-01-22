package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// NotificationService handles sending notifications via Apprise
type NotificationService struct {
	appriseURL string
	httpClient *http.Client
}

// apprisePayload represents the JSON payload for Apprise API
type apprisePayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Type  string `json:"type"`
}

// NewEmailService creates a new NotificationService instance
// Note: Function name kept for backward compatibility with existing code
func NewEmailService(appriseURL string, _ string) *NotificationService {
	if appriseURL == "" {
		appriseURL = os.Getenv("APPRISE_URL")
	}

	if appriseURL == "" {
		panic("APPRISE_URL environment variable is not set")
	}

	return &NotificationService{
		appriseURL: appriseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendEmail sends a notification via Apprise
// Note: Function name kept for backward compatibility with existing code
func (s *NotificationService) SendEmail(ctx context.Context, subject, body string) error {
	// Determine notification type based on subject content
	notificationType := "info"
	subjectLower := strings.ToLower(subject)
	if strings.Contains(subjectLower, "error") || strings.Contains(subjectLower, "failed") {
		notificationType = "failure"
	}

	payload := apprisePayload{
		Title: subject,
		Body:  body,
		Type:  notificationType,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.appriseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create notification request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("apprise returned non-success status: %d", resp.StatusCode)
	}

	return nil
}
