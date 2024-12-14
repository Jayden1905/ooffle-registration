package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"

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
	verificationLink := fmt.Sprintf("%s/api/v1/user/verify/email?token=%s", config.Envs.BackendHost, token)

	// Load the HTML template
	tmplPath := "templates/verify_email.html"
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		log.Printf("Error reading email template: %v", err)
		return err
	}

	// Parse the template
	tmpl, err := template.New("verify_email").Parse(string(tmplContent))
	if err != nil {
		log.Printf("Error parsing email template: %v", err)
		return err
	}

	// Prepare template data
	data := struct {
		VerificationLink string
	}{
		VerificationLink: verificationLink,
	}

	// Render the template
	var renderedBody bytes.Buffer
	if err := tmpl.Execute(&renderedBody, data); err != nil {
		log.Printf("Error executing email template: %v", err)
		return err
	}

	// Create the email content
	subject := "Subject: Verify Your Account\r\n"
	contentType := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg := []byte(subject + contentType + "\r\n" + renderedBody.String())

	// Send the email
	err = smtp.SendMail(es.SMTPHost+":"+es.SMTPPort, auth, es.FromEmail, []string{toEmail}, msg)
	if err != nil {
		log.Printf("Error sending email to %s: %v", toEmail, err)
		return err
	}

	return nil
}
