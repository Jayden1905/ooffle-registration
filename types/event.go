package types

import (
	"context"
	"time"
)

type Event struct {
	EventID     int32     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location"`
	UserID      int32     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventStore interface {
	CreateNewEvent(ctx context.Context, event *Event) error
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, eventID int32) error
	DeleteAllEvents(ctx context.Context, userID int32) error
	GetAllEvents(userID int32) ([]*Event, error)
	GetEventByTitle(title string) (*Event, error)
	GetEventByID(eventID int32) (*Event, error)
}

type CreateEventPayload struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	Location    string    `json:"location" validate:"required"`
}
