package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/jayden1905/event-registration-software/config"
	"github.com/jayden1905/event-registration-software/types"
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

// SendInvitationEmail sends the invitation email
func (es *EmailService) SendInvitationEmail(attendee *types.Attendee, template *types.EmailTemplate) error {
	auth := smtp.PlainAuth("", es.SMTPUsername, es.SMTPPassword, es.SMTPHost)

	// Replace template variables with attendee data
	content := template.Content
	content = strings.Replace(content, "{{first_name}}", attendee.FirstName, -1)
	content = strings.Replace(content, "{{last_name}}", attendee.LastName, -1)
	content = strings.Replace(content, "{{qr_code}}", attendee.QrCode, -1) // Use cid for the inline image

	// Prepare the email body
	body := fmt.Sprintf(`
		%s
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Email Preview</title>
				<style>
					body {
						font-family: Arial, sans-serif;
						line-height: 1.6;
						margin: 0;
						padding: 0;
						color: black;
					}
					.container {
						max-width: 600px;
						width: 100%%;
						margin: 0 auto;
					}
					.img-container img {
						width: 100%%;
						height: auto;
						object-fit: cover;
						object-position: center;
					}
				</style>
			</head>
		<body>
		<table style="background-color: %s;" class="container" role="presentation" cellspacing="0" cellpadding="0">
			<tr>
				<td class="img-container">
					<img src="%s" alt="Header" />
				</td>
			</tr>
			<tr>
				<td>
					<div>%s</div>
				</td>
			</tr>
			<tr>
				<td class="img-container">
					<img src="%s" alt="Footer" />
				</td>
			</tr>
		</table>
		</body>
		</html>
	`, template.Message, template.BgColor, template.HeaderImage, content, template.FooterImage)

	subject := fmt.Sprintf("Subject: %s\r\n", template.Subject)
	contentType := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n"
	msg := []byte(subject + contentType + "\r\n" + body)

	// Send the email
	err := smtp.SendMail(es.SMTPHost+":"+es.SMTPPort, auth, es.FromEmail, []string{attendee.Email}, msg)
	if err != nil {
		log.Printf("Error sending email to %s: %v", attendee.Email, err)
		return err
	}

	log.Printf("Invitation email sent to %s", attendee.Email)

	return nil
}
