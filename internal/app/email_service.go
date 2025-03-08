package app

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     os.Getenv("SMTP_PORT"),
		smtpUsername: os.Getenv("SMTP_USERNAME"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
		fromEmail:    os.Getenv("SMTP_FROM_EMAIL"),
	}
}

func (s *EmailService) SendOTP(ctx context.Context, email, otp string) error {
	// Email template
	subject := "Your Email Verification Code"
	body := fmt.Sprintf(`
		Hello,

		Your email verification code is: %s

		This code will expire in 15 minutes.

		If you didn't request this code, please ignore this email.

		Best regards,
		Cookaholic Team
	`, otp)

	// Prepare email message
	message := fmt.Sprintf("Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", subject, body)

	// Connect to SMTP server
	var auth smtp.Auth
	if s.smtpUsername != "" && s.smtpPassword != "" {
		auth = smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	}

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	// Send email
	return smtp.SendMail(addr, auth, s.fromEmail, []string{email}, []byte(message))
}
