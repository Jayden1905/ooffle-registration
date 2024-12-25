package email

import "github.com/jayden1905/event-registration-software/types"

type Mailer interface {
	SendVerificationEmail(toEmail string, token string) error
	SendInvitationEmail(attendee *types.Attendee, template *types.EmailTemplate) error
}
