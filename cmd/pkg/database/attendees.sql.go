// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: attendees.sql

package database

import (
	"context"
	"database/sql"
)

const createAttendee = `-- name: CreateAttendee :exec
INSERT INTO attendees (first_name, last_name, email, qr_code, company_name, title, table_no, role, attendence, event_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type CreateAttendeeParams struct {
	FirstName   string
	LastName    string
	Email       string
	QrCode      sql.NullString
	CompanyName sql.NullString
	Title       sql.NullString
	TableNo     sql.NullInt32
	Role        sql.NullString
	Attendence  NullAttendeesAttendence
	EventID     int32
}

func (q *Queries) CreateAttendee(ctx context.Context, arg CreateAttendeeParams) error {
	_, err := q.db.ExecContext(ctx, createAttendee,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.QrCode,
		arg.CompanyName,
		arg.Title,
		arg.TableNo,
		arg.Role,
		arg.Attendence,
		arg.EventID,
	)
	return err
}

const deleteAllAttendeesByEventID = `-- name: DeleteAllAttendeesByEventID :exec
DELETE FROM attendees WHERE event_id = ?
`

func (q *Queries) DeleteAllAttendeesByEventID(ctx context.Context, eventID int32) error {
	_, err := q.db.ExecContext(ctx, deleteAllAttendeesByEventID, eventID)
	return err
}

const deleteAttendeeByID = `-- name: DeleteAttendeeByID :exec
DELETE FROM attendees WHERE id = ?
`

func (q *Queries) DeleteAttendeeByID(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteAttendeeByID, id)
	return err
}

const getAllAttendeesPaginatedByEventID = `-- name: GetAllAttendeesPaginatedByEventID :many
SELECT id, first_name, last_name, email, qr_code, company_name, title, table_no, role, attendence, event_id FROM attendees WHERE event_id = ? LIMIT ? OFFSET ?
`

type GetAllAttendeesPaginatedByEventIDParams struct {
	EventID int32
	Limit   int32
	Offset  int32
}

func (q *Queries) GetAllAttendeesPaginatedByEventID(ctx context.Context, arg GetAllAttendeesPaginatedByEventIDParams) ([]Attendee, error) {
	rows, err := q.db.QueryContext(ctx, getAllAttendeesPaginatedByEventID, arg.EventID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Attendee
	for rows.Next() {
		var i Attendee
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.QrCode,
			&i.CompanyName,
			&i.Title,
			&i.TableNo,
			&i.Role,
			&i.Attendence,
			&i.EventID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAttendeeByEmail = `-- name: GetAttendeeByEmail :one
SELECT id, first_name, last_name, email, qr_code, company_name, title, table_no, role, attendence, event_id FROM attendees WHERE email = ?
`

func (q *Queries) GetAttendeeByEmail(ctx context.Context, email string) (Attendee, error) {
	row := q.db.QueryRowContext(ctx, getAttendeeByEmail, email)
	var i Attendee
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.QrCode,
		&i.CompanyName,
		&i.Title,
		&i.TableNo,
		&i.Role,
		&i.Attendence,
		&i.EventID,
	)
	return i, err
}

const getAttendeeByID = `-- name: GetAttendeeByID :one
SELECT id, first_name, last_name, email, qr_code, company_name, title, table_no, role, attendence, event_id FROM attendees WHERE id = ?
`

func (q *Queries) GetAttendeeByID(ctx context.Context, id int32) (Attendee, error) {
	row := q.db.QueryRowContext(ctx, getAttendeeByID, id)
	var i Attendee
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.QrCode,
		&i.CompanyName,
		&i.Title,
		&i.TableNo,
		&i.Role,
		&i.Attendence,
		&i.EventID,
	)
	return i, err
}

const saveAttendeeWithQRCode = `-- name: SaveAttendeeWithQRCode :exec
UPDATE attendees SET qr_code = ? WHERE id = ?
`

type SaveAttendeeWithQRCodeParams struct {
	QrCode sql.NullString
	ID     int32
}

func (q *Queries) SaveAttendeeWithQRCode(ctx context.Context, arg SaveAttendeeWithQRCodeParams) error {
	_, err := q.db.ExecContext(ctx, saveAttendeeWithQRCode, arg.QrCode, arg.ID)
	return err
}
