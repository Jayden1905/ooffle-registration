package event

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

// GetAllEvents fetches all events from the database
func (s *Store) GetAllEvents(userID int32) ([]*types.Event, error) {
	events, err := s.db.GetAllEventsByUserID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// Convert the database event to the event type
	var allEvents []*types.Event

	for _, event := range events {
		allEvents = append(allEvents, &types.Event{
			EventID:     int32(event.EventID),
			Title:       event.Title,
			Description: event.Description,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			Location:    event.Location,
			UserID:      int32(event.UserID),
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		})
	}

	return allEvents, nil
}

// CreateEvent creates a new event in the database
func (s *Store) CreateNewEvent(ctx context.Context, event *types.Event) error {
	err := s.db.CreateEvent(ctx, database.CreateEventParams{
		Title:       event.Title,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
		UserID:      event.UserID,
	})
	if err != nil {
		return err
	}

	return nil
}

// UpdateEvent updates an existing event in the database
func (s *Store) UpdateEvent(ctx context.Context, event *types.Event) error {
	err := s.db.UpdateEventByID(ctx, database.UpdateEventByIDParams{
		EventID:     event.EventID,
		Title:       event.Title,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteEvent deletes an event from the database
func (s *Store) DeleteEvent(ctx context.Context, eventID int32) error {
	err := s.db.DeleteEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAllEvents deletes all events from the database
func (s *Store) DeleteAllEvents(ctx context.Context, userID int32) error {
	err := s.db.DeleteAllEventsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

// GetEventByTitle fetches an event by its title
func (s *Store) GetEventByTitle(title string) (*types.Event, error) {
	event, err := s.db.GetEventByTitle(context.Background(), title)
	if err != nil {
		return nil, err
	}

	return &types.Event{
		EventID:     int32(event.EventID),
		Title:       event.Title,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
		UserID:      int32(event.UserID),
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}, nil
}

// GetEventByID fetches an event by its ID
func (s *Store) GetEventByID(eventID int32) (*types.Event, error) {
	event, err := s.db.GetEventByID(context.Background(), eventID)
	if err != nil {
		return nil, err
	}

	return &types.Event{
		EventID:     int32(event.EventID),
		Title:       event.Title,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
		UserID:      int32(event.UserID),
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}, nil
}
