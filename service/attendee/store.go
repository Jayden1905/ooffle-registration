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
func (s *Store) CreateAttendee(ctx context.Context, attendee *types.Attendee) error {
	attendanceValue := database.AttendeesAttendanceNo
	if attendee.Attendance {
		attendanceValue = database.AttendeesAttendanceYes
	}

	err := s.db.CreateAttendee(ctx, database.CreateAttendeeParams{
		FirstName:   attendee.FristName,
		LastName:    attendee.LastName,
		Email:       attendee.Email,
		EventID:     attendee.EventID,
		QrCode:      sql.NullString{String: attendee.QrCode, Valid: true},
		CompanyName: sql.NullString{String: attendee.CompanyName, Valid: true},
		Title:       sql.NullString{String: attendee.Title, Valid: true},
		TableNo:     sql.NullInt32{Int32: attendee.TableNo, Valid: true},
		Role:        sql.NullString{String: attendee.Role, Valid: true},
		Attendance: database.NullAttendeesAttendance{
			AttendeesAttendance: attendanceValue,
			Valid:               true,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// GetAttendeeByEmail fetches an attendee from the database by email
func (s *Store) GetAttendeeByEmail(email string) (*types.Attendee, error) {
	attendee, err := s.db.GetAttendeeByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	return &types.Attendee{
		ID:          attendee.ID,
		FristName:   attendee.FirstName,
		LastName:    attendee.LastName,
		Email:       attendee.Email,
		EventID:     attendee.EventID,
		QrCode:      attendee.QrCode.String,
		CompanyName: attendee.CompanyName.String,
		Title:       attendee.Title.String,
		TableNo:     attendee.TableNo.Int32,
		Role:        attendee.Role.String,
		Attendance:  attendee.Attendance.Valid,
	}, nil
}

func (s *Store) GetAttendeeByID(id int32) (*types.Attendee, error) {
	attendee, err := s.db.GetAttendeeByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &types.Attendee{
		ID:          attendee.ID,
		FristName:   attendee.FirstName,
		LastName:    attendee.LastName,
		Email:       attendee.Email,
		EventID:     attendee.EventID,
		QrCode:      attendee.QrCode.String,
		CompanyName: attendee.CompanyName.String,
		Title:       attendee.Title.String,
		TableNo:     attendee.TableNo.Int32,
		Role:        attendee.Role.String,
		Attendance:  attendee.Attendance.Valid,
	}, nil
}

// DeleteAttendeeByID deletes an attendee from the database by ID
func (s *Store) DeleteAttendeeByID(attendeeID int32) error {
	err := s.db.DeleteAttendeeByID(context.Background(), attendeeID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAllAttendeesByEventID deletes all attendees from the database by event ID
func (s *Store) DeleteAllAttendeesByEventID(eventID int32) error {
	err := s.db.DeleteAllAttendeesByEventID(context.Background(), eventID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateAttendeeByID updates an attendee in the database by ID
func (s *Store) UpdateAttendeeByID(attendeeID int32, data *types.Attendee) error {
	attendanceValue := database.AttendeesAttendanceNo
	if data.Attendance {
		attendanceValue = database.AttendeesAttendanceYes
	}

	err := s.db.UpdateAttendeeByID(context.Background(), database.UpdateAttendeeByIDParams{
		ID:          attendeeID,
		FirstName:   data.FristName,
		LastName:    data.LastName,
		Email:       data.Email,
		QrCode:      sql.NullString{String: data.QrCode, Valid: true},
		CompanyName: sql.NullString{String: data.CompanyName, Valid: true},
		Title:       sql.NullString{String: data.Title, Valid: true},
		TableNo:     sql.NullInt32{Int32: data.TableNo, Valid: true},
		Role:        sql.NullString{String: data.Role, Valid: true},
		Attendance: database.NullAttendeesAttendance{
			AttendeesAttendance: attendanceValue,
			Valid:               true,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// GetAllAttendeesPaginated fetches all attendees from the database with pagination
func (s *Store) GetAllAttendeesPaginated(page int32, pageSize int32, eventID int32) ([]*types.Attendee, error) {
	offset := (page - 1) * pageSize

	attendees, err := s.db.GetAllAttendeesPaginatedByEventID(context.Background(), database.GetAllAttendeesPaginatedByEventIDParams{
		Limit:   pageSize,
		Offset:  offset,
		EventID: eventID,
	})
	if err != nil {
		return nil, err
	}

	var allAttendees []*types.Attendee

	for _, attendee := range attendees {
		allAttendees = append(allAttendees, &types.Attendee{
			ID:          attendee.ID,
			FristName:   attendee.FirstName,
			LastName:    attendee.LastName,
			Email:       attendee.Email,
			EventID:     attendee.EventID,
			QrCode:      attendee.QrCode.String,
			CompanyName: attendee.CompanyName.String,
			Title:       attendee.Title.String,
			TableNo:     attendee.TableNo.Int32,
			Role:        attendee.Role.String,
			Attendance:  attendee.Attendance.Valid,
		})
	}

	return allAttendees, nil
}
