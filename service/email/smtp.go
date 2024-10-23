package email

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/jayden1905/event-registration-software/config"
)

// EmailService holds the SMTP server information for sending emails
type EmailService struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

// NewEmailService creates a new EmailService instance
func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost:     config.Envs.SMPTHost,
		SMTPPort:     config.Envs.SMTPPort,
		SMTPUsername: config.Envs.SMTPUsername,
		SMTPPassword: config.Envs.SMTPPassword,
		FromEmail:    "noreply@yourdomain.com",
	}
}

// SendVerificationEmail sends a verification email with a token link in HTML format
func (es *EmailService) SendVerificationEmail(toEmail string, token string) error {
	auth := smtp.PlainAuth("", es.SMTPUsername, es.SMTPPassword, es.SMTPHost)

	// Verification link
	verificationLink := fmt.Sprintf("%s/api/v1/user/verify/email?token=%s", "http://127.0.0.1:8080", token)

	// Create the HTML email body
	subject := "Subject: Verify Your Account\r\n"
	contentType := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	body := fmt.Sprintf(`<html><body><p>Please verify your account by clicking the link below:</p><a href="%s">Verify Your Account</a></body></html>`, verificationLink)
	msg := []byte(subject + contentType + "\r\n" + body)

	// Send the email
	err := smtp.SendMail(es.SMTPHost+":"+es.SMTPPort, auth, es.FromEmail, []string{toEmail}, msg)
	if err != nil {
		log.Printf("Error sending email to %s: %v", toEmail, err)
		return err
	}

	return nil
}
