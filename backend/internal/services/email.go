package services

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// EmailService handles operations related to sending emails
type EmailService struct {
	recipient string
	from      string
	username  string
	password  string
	template  string
}

// NewEmailService creates a new EmailService instance
func NewEmailService(recipient string, templatePath string) *EmailService {
	from := os.Getenv("SMTP2GO_FROM_EMAIL")
	username := os.Getenv("SMTP2GO_USERNAME")
	password := os.Getenv("SMTP2GO_PASSWORD")

	// Validate required fields
	if from == "" {
		panic("SMTP2GO_FROM_EMAIL environment variable is not set")
	}
	if username == "" {
		panic("SMTP2GO_USERNAME environment variable is not set")
	}
	if password == "" {
		panic("SMTP2GO_PASSWORD environment variable is not set")
	}
	if recipient == "" {
		panic("recipient email address is required")
	}

	// Basic email format validation
	if !strings.Contains(from, "@") {
		panic("SMTP2GO_FROM_EMAIL is not a valid email address")
	}
	if !strings.Contains(recipient, "@") {
		panic("recipient is not a valid email address")
	}

	// Load email template if provided
	var template string
	if templatePath != "" {
		templateBytes, err := os.ReadFile(templatePath)
		if err != nil {
			panic(fmt.Sprintf("failed to read email template: %v", err))
		}
		template = string(templateBytes)

		// Replace dashboard URL placeholder if environment variable is set
		dashboardURL := os.Getenv("DASHBOARD_URL")
		if dashboardURL != "" {
			template = strings.Replace(template, "{{DASHBOARD_URL}}", dashboardURL, -1)
		}
	}

	return &EmailService{
		recipient: recipient,
		from:      from,
		username:  username,
		password:  password,
		template:  template,
	}
}

// SendEmail sends an email to the specified address with the given subject and body
func (s *EmailService) SendEmail(ctx context.Context, subject, body string) error {
	// SMTP2Go settings
	host := "mail.smtp2go.com"
	port := "587"
	auth := smtp.PlainAuth("", s.username, s.password, host)

	// Use template if available, otherwise use plain text
	content := body
	if s.template != "" {
		content = fmt.Sprintf(s.template, body)
	}

	// Format the email with HTML content type
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", s.from, s.recipient, subject, content)

	// Send the email
	err := smtp.SendMail(
		host+":"+port,
		auth,
		s.from,
		[]string{s.recipient},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
