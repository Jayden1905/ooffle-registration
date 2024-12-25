package types

import "context"

type EmailTemplate struct {
	ID          int32  `json:"id"`
	EventID     int32  `json:"event_id"`
	HeaderImage string `json:"header_image"`
	Content     string `json:"content"`
	FooterImage string `json:"footer_image"`
}

type EmailTempalteStore interface {
	GetEmailTemplateByEventID(ctx context.Context, eventID int32) (*EmailTemplate, error)
	CreateEmailTemplate(ctx context.Context, emailTemplate *EmailTemplate) error
	UpdateEmailTemplate(ctx context.Context, emailTemplate *EmailTemplate) error
}

type CreateEmailTemplatePayload struct {
	EventID     int32  `json:"event_id" validate:"required"`
	HeaderImage string `json:"header_image" validate:"required"`
	Content     string `json:"content" validate:"required"`
	FooterImage string `json:"footer_image" validate:"required"`
}

type UpdateEmailTemplatePayload struct {
	ID          int32  `json:"id"`
	EventID     int32  `json:"event_id" validate:"required"`
	HeaderImage string `json:"header_image" validate:"required"`
	Content     string `json:"content" validate:"required"`
	FooterImage string `json:"footer_image" validate:"required"`
}
