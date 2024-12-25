// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type AttendeesAttendance string

const (
	AttendeesAttendanceYes AttendeesAttendance = "Yes"
	AttendeesAttendanceNo  AttendeesAttendance = "No"
)

func (e *AttendeesAttendance) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AttendeesAttendance(s)
	case string:
		*e = AttendeesAttendance(s)
	default:
		return fmt.Errorf("unsupported scan type for AttendeesAttendance: %T", src)
	}
	return nil
}

type NullAttendeesAttendance struct {
	AttendeesAttendance AttendeesAttendance
	Valid               bool // Valid is true if AttendeesAttendance is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAttendeesAttendance) Scan(value interface{}) error {
	if value == nil {
		ns.AttendeesAttendance, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AttendeesAttendance.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAttendeesAttendance) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AttendeesAttendance), nil
}

type RolesName string

const (
	RolesNameSuperUser  RolesName = "super_user"
	RolesNameNormalUser RolesName = "normal_user"
)

func (e *RolesName) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RolesName(s)
	case string:
		*e = RolesName(s)
	default:
		return fmt.Errorf("unsupported scan type for RolesName: %T", src)
	}
	return nil
}

type NullRolesName struct {
	RolesName RolesName
	Valid     bool // Valid is true if RolesName is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRolesName) Scan(value interface{}) error {
	if value == nil {
		ns.RolesName, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RolesName.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRolesName) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RolesName), nil
}

type SubscriptionsStatus string

const (
	SubscriptionsStatusActive    SubscriptionsStatus = "Active"
	SubscriptionsStatusInactive  SubscriptionsStatus = "Inactive"
	SubscriptionsStatusPending   SubscriptionsStatus = "Pending"
	SubscriptionsStatusCancelled SubscriptionsStatus = "Cancelled"
)

func (e *SubscriptionsStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = SubscriptionsStatus(s)
	case string:
		*e = SubscriptionsStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for SubscriptionsStatus: %T", src)
	}
	return nil
}

type NullSubscriptionsStatus struct {
	SubscriptionsStatus SubscriptionsStatus
	Valid               bool // Valid is true if SubscriptionsStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullSubscriptionsStatus) Scan(value interface{}) error {
	if value == nil {
		ns.SubscriptionsStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.SubscriptionsStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullSubscriptionsStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.SubscriptionsStatus), nil
}

type Attendee struct {
	ID          int32
	FirstName   string
	LastName    string
	Email       string
	QrCode      sql.NullString
	CompanyName sql.NullString
	Title       sql.NullString
	TableNo     sql.NullInt32
	Role        sql.NullString
	Attendance  NullAttendeesAttendance
	EventID     int32
}

type AttendeesCustomField struct {
	ID         int32
	AttendeeID int32
	FieldName  sql.NullString
	FieldValue sql.NullString
	FieldType  sql.NullString
}

type EmailTemplate struct {
	ID          int32
	EventID     int32
	HeaderImage sql.NullString
	Content     sql.NullString
	FooterImage sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Subject     sql.NullString
	Message     sql.NullString
	BgColor     sql.NullString
}

type Event struct {
	EventID     int32
	Title       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Location    string
	UserID      int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Role struct {
	RoleID int8
	Name   RolesName
}

type Subscription struct {
	SubscriptionID int8
	Status         SubscriptionsStatus
}

type User struct {
	UserID         int32
	RoleID         int8
	FirstName      string
	LastName       string
	Email          string
	Password       string
	SubscriptionID int8
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Verify         bool
}
