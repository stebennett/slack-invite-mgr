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
	tag        string
	httpClient *http.Client
}

// apprisePayload represents the JSON payload for Apprise API
type apprisePayload struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Type   string `json:"type"`
	Format string `json:"format"`
	Tag    string `json:"tag,omitempty"`
}

// NewEmailService creates a new NotificationService instance
// Note: Function name kept for backward compatibility with existing code
func NewEmailService(appriseURL string, tag string) *NotificationService {
	if appriseURL == "" {
		appriseURL = os.Getenv("APPRISE_URL")
	}

	if appriseURL == "" {
		panic("APPRISE_URL environment variable is not set")
	}

	if tag == "" {
		tag = os.Getenv("APPRISE_TAG")
	}

	return &NotificationService{
		appriseURL: appriseURL,
		tag:        tag,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendEmail sends an HTML notification via Apprise
// Note: Function name kept for backward compatibility with existing code
func (s *NotificationService) SendEmail(ctx context.Context, subject, body string) error {
	// Determine notification type based on subject content
	notificationType := "info"
	subjectLower := strings.ToLower(subject)
	if strings.Contains(subjectLower, "error") || strings.Contains(subjectLower, "failed") {
		notificationType = "failure"
	}

	// Wrap body in HTML template
	htmlBody := formatHTMLEmail(subject, body, notificationType)

	payload := apprisePayload{
		Title:  subject,
		Body:   htmlBody,
		Type:   notificationType,
		Format: "html",
		Tag:    s.tag,
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

// formatHTMLEmail wraps the message body in an HTML template
func formatHTMLEmail(subject, body, notificationType string) string {
	// Set colors based on notification type
	headerColor := "#4A154B" // Slack purple (default/info)
	if notificationType == "failure" {
		headerColor = "#E01E5A" // Red for errors
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333333; margin: 0; padding: 0; background-color: #f4f4f4;">
    <div style="max-width: 600px; margin: 20px auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);">
        <div style="background-color: %s; color: #ffffff; padding: 20px; border-radius: 8px 8px 0 0; text-align: center;">
            <h1 style="margin: 0; font-size: 24px; font-weight: 600;">Slack Invite Manager</h1>
        </div>
        <div style="padding: 30px; background-color: #ffffff;">
            <div style="font-size: 16px; color: #333333; margin-bottom: 20px;">
                %s
            </div>
        </div>
        <div style="padding: 20px; text-align: center; font-size: 14px; color: #666666; border-top: 1px solid #eeeeee; background-color: #fafafa; border-radius: 0 0 8px 8px;">
            <p style="margin: 0;">This is an automated message from the Slack Invite Manager system.</p>
        </div>
    </div>
</body>
</html>`, subject, headerColor, body)
}
