package attendee

import (
	"context"
	"database/sql"

	"github.com/jayden1905/event-registration-software/cmd/pkg/database"
	"github.com/jayden1905/event-registration-software/types"
)

type Store struct {
	db *database.Queries
}

func NewStore(db *database.Queries) *Store {
	return &Store{db: db}
}

// CreateAttendee creates a new attendee in the database
func (s *Store) CreateAttendee(attendee *types.Attendee) error {
	err := s.db.CreateAttendee(context.Background(), database.CreateAttendeeParams{
		FirstName:   attendee.FristName,
		LastName:    attendee.LastName,
		Email:       attendee.Email,
		EventID:     attendee.EventID,
		QrCode:      sql.NullString{String: attendee.QrCode, Valid: true},
		CompanyName: sql.NullString{String: attendee.CompanyName, Valid: true},
		Title:       sql.NullString{String: attendee.Title, Valid: true},
		TableNo:     sql.NullInt32{Int32: attendee.TableNo, Valid: true},
		Role:        sql.NullString{String: attendee.Role, Valid: true},
		Attendence:  database.NullAttendeesAttendence{Valid: attendee.Attendence},
	})
	if err != nil {
		return err
	}

	return nil
}

// GetAllAttendeesPaginated fetches all attendees from the database with pagination
func (s *Store) GetAllAttendeesPaginated(page int32, pageSize int32) ([]*types.Attendee, error) {
	offset := (page - 1) * pageSize

	attendees, err := s.db.GetAllAttendeesPaginated(context.Background(), database.GetAllAttendeesPaginatedParams{
		Limit:  pageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var allAttendees []*types.Attendee

	for _, attendee := range attendees {
		allAttendees = append(allAttendees, &types.Attendee{
			FristName:   attendee.FirstName,
			LastName:    attendee.LastName,
			Email:       attendee.Email,
			EventID:     attendee.EventID,
			QrCode:      attendee.QrCode.String,
			CompanyName: attendee.CompanyName.String,
			Title:       attendee.Title.String,
			TableNo:     attendee.TableNo.Int32,
			Role:        attendee.Role.String,
			Attendence:  attendee.Attendence.Valid,
		})
	}

	return allAttendees, nil
}
