-- name: CreateAttendee :exec
INSERT INTO attendees (first_name, last_name, email, qr_code, company_name, title, table_no, role, attendence, event_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: FindAttendeeByEmail :one
SELECT * FROM attendees WHERE email = ?;

-- name: GetAllAttendees :many
SELECT * FROM attendees;

-- name: GetAllAttendeesPaginated :many
SELECT * FROM attendees LIMIT ? OFFSET ?;

-- name: DeleteAttendeeByID :exec
DELETE FROM attendees WHERE id = ?;

-- name: DeleteAllAttendees :exec
DELETE FROM attendees;

-- name: SaveAttendeeWithQRCode :exec
UPDATE attendees SET qr_code = ? WHERE id = ?;
