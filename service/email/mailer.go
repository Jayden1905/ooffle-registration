package email

type Mailer interface {
	SendVerificationEmail(toEmail string, token string) error
}
