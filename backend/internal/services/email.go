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
}

// NewEmailService creates a new EmailService instance
func NewEmailService(recipient string) *EmailService {
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

	return &EmailService{
		recipient: recipient,
		from:      from,
		username:  username,
		password:  password,
	}
}

// SendEmail sends an email to the specified address with the given subject and body
func (s *EmailService) SendEmail(ctx context.Context, subject, body string) error {
	// SMTP2Go settings
	host := "mail.smtp2go.com"
	port := "587"
	auth := smtp.PlainAuth("", s.username, s.password, host)

	// Format the email
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.from, s.recipient, subject, body)

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
