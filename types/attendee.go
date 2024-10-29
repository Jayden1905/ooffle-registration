package types

type Attendee struct {
	FristName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	EventID     int32  `json:"event_id"`
	QrCode      string `json:"qr_code"`
	CompanyName string `json:"company_name"`
	Title       string `json:"title"`
	TableNo     int32  `json:"table_no"`
	Role        string `json:"role"`
	Attendence  bool   `json:"attendence"`
}

type AttendeeStore interface {
	GetAllAttendeesPaginated(page int32, pageSize int32) ([]*Attendee, error)
	CreateAttendee(attendee *Attendee) error
}
