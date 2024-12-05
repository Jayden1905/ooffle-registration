package types

import "context"

type Attendee struct {
	ID          int32  `json:"id"`
	FristName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	EventID     int32  `json:"event_id"`
	QrCode      string `json:"qr_code"`
	CompanyName string `json:"company_name"`
	Title       string `json:"title"`
	TableNo     int32  `json:"table_no"`
	Role        string `json:"role"`
	Attendance  bool   `json:"attendance"`
}

type AttendeeStore interface {
	GetAllAttendeesPaginated(page int32, pageSize int32, eventID int32) ([]*Attendee, error)
	GetAttendeeByEmail(email string) (*Attendee, error)
	GetAttendeeByID(attendeeID int32) (*Attendee, error)
	CreateAttendee(ctx context.Context, attendee *Attendee) error
	DeleteAttendeeByID(attendeeID int32) error
	DeleteAllAttendeesByEventID(eventID int32) error
	UpdateAttendeeByID(attendeeID int32, data *Attendee) error
}

type CreateAttendeePayload struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	EventID     int32  `json:"event_id" validate:"required"`
	CompanyName string `json:"company_name"`
	Title       string `json:"title"`
	TableNo     int32  `json:"table_no"`
	Role        string `json:"role"`
}

type UpdateAttendeePayload struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	Title       string `json:"title"`
	TableNo     int32  `json:"table_no"`
	Role        string `json:"role"`
	Attendance  bool   `json:"attendance"`
}
