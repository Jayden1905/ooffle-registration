package email

import (
	"context"

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
		HeaderImage: emailTemplate.HeaderImage,
		Content:     emailTemplate.Content,
		FooterImage: emailTemplate.FooterImage,
	}, nil
}

// CreateEmailTemplate creates an email template in the database
func (s *Store) CreateEmailTemplate(ctx context.Context, emailTemplate *types.EmailTemplate) error {
	err := s.db.CreateEmailTemplate(ctx, database.CreateEmailTemplateParams{
		EventID:     emailTemplate.EventID,
		HeaderImage: emailTemplate.HeaderImage,
		Content:     emailTemplate.Content,
		FooterImage: emailTemplate.FooterImage,
	})

	if err != nil {
		return err
	}

	return nil
}

// UpdateEmailTemplate updates an email template in the database
func (s *Store) UpdateEmailTemplate(ctx context.Context, emailTemplate *types.EmailTemplate) error {
	err := s.db.UpdateEmailTemplateByID(ctx, database.UpdateEmailTemplateByIDParams{
		EventID:     emailTemplate.EventID,
		HeaderImage: emailTemplate.HeaderImage,
		Content:     emailTemplate.Content,
		FooterImage: emailTemplate.FooterImage,
		ID:          emailTemplate.ID,
	})

	if err != nil {
		return err
	}

	return nil
}
