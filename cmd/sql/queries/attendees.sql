-- name: CreateAttendee :exec
INSERT INTO attendees (first_name, last_name, email, qr_code, company_name, title, table_no, role, attendence, event_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetAttendeeByEmail :one
SELECT * FROM attendees WHERE email = ?;

-- name: GetAttendeeByID :one
SELECT * FROM attendees WHERE id = ?;

-- name: GetAllAttendeesPaginatedByEventID :many
SELECT * FROM attendees WHERE event_id = ? LIMIT ? OFFSET ?;

-- name: DeleteAttendeeByID :exec
DELETE FROM attendees WHERE id = ?;

-- name: DeleteAllAttendeesByEventID :exec
DELETE FROM attendees WHERE event_id = ?;

-- name: SaveAttendeeWithQRCode :exec
UPDATE attendees SET qr_code = ? WHERE id = ?;
