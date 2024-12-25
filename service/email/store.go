package email

import (
	"context"
	"database/sql"

	"github.com/jayden1905/event-registration-software/cmd/pkg/database"
	"github.com/jayden1905/event-registration-software/types"
)

type Store struct {
	db *database.Queries
}

// NewStore initializes the Store with the database queries
func NewStore(db *database.Queries) *Store {
	return &Store{db: db}
}

// GetEmailTemplateByEventID fetches an email template by its event ID from the database
func (s *Store) GetEmailTemplateByEventID(c context.Context, eventID int32) (*types.EmailTemplate, error) {
	emailTemplate, err := s.db.GetEmailTemplateByEventID(c, eventID)
	if err != nil {
		return nil, err
	}

	return &types.EmailTemplate{
		ID:          int32(emailTemplate.ID),
		EventID:     int32(emailTemplate.EventID),
		HeaderImage: emailTemplate.HeaderImage.String,
		Content:     emailTemplate.Content.String,
		FooterImage: emailTemplate.FooterImage.String,
		Subject:     emailTemplate.Subject.String,
		BgColor:     emailTemplate.BgColor.String,
		Message:     emailTemplate.Message.String,
	}, nil
}

// CreateEmailTemplate creates an email template in the database
func (s *Store) CreateEmailTemplate(ctx context.Context, emailTemplate *types.EmailTemplate) error {
	err := s.db.CreateEmailTemplate(ctx, database.CreateEmailTemplateParams{
		EventID:     emailTemplate.EventID,
		HeaderImage: sql.NullString{String: emailTemplate.HeaderImage, Valid: emailTemplate.HeaderImage != ""},
		Content:     sql.NullString{String: emailTemplate.Content, Valid: emailTemplate.Content != ""},
		FooterImage: sql.NullString{String: emailTemplate.FooterImage, Valid: emailTemplate.FooterImage != ""},
		Subject:     sql.NullString{String: emailTemplate.Subject, Valid: emailTemplate.Subject != ""},
		BgColor:     sql.NullString{String: emailTemplate.BgColor, Valid: emailTemplate.BgColor != ""},
		Message:     sql.NullString{String: emailTemplate.Message, Valid: emailTemplate.Message != ""},
	})

	if err != nil {
		return err
	}

	return nil
}

// UpdateEmailTemplate updates an email template in the database
func (s *Store) UpdateEmailTemplate(ctx context.Context, emailTemplate *types.EmailTemplate) error {
	err := s.db.UpdateEmailTemplateByID(ctx, database.UpdateEmailTemplateByIDParams{
		ID:          emailTemplate.ID,
		EventID:     emailTemplate.EventID,
		HeaderImage: sql.NullString{String: emailTemplate.HeaderImage, Valid: emailTemplate.HeaderImage != ""},
		Content:     sql.NullString{String: emailTemplate.Content, Valid: emailTemplate.Content != ""},
		FooterImage: sql.NullString{String: emailTemplate.FooterImage, Valid: emailTemplate.FooterImage != ""},
		Subject:     sql.NullString{String: emailTemplate.Subject, Valid: emailTemplate.Subject != ""},
		BgColor:     sql.NullString{String: emailTemplate.BgColor, Valid: emailTemplate.BgColor != ""},
		Message:     sql.NullString{String: emailTemplate.Message, Valid: emailTemplate.Message != ""},
	})

	if err != nil {
		return err
	}

	return nil
}
