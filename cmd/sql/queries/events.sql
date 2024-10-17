-- name: CreateEvent :exec
INSERT INTO events (
        title,
        description,
        start_date,
        end_date,
        location,
        user_id
    )
VALUES (?, ?, ?, ?, ?, ?);
-- name: UpdateEventByID :exec
UPDATE events
SET title = ?,
    description = ?,
    start_date = ?,
    end_date = ?,
    location = ?
WHERE event_id = ?;
-- name: DeleteEventByID :exec
DELETE FROM events
WHERE event_id = ?;
-- name: DeleteAllEventsByUserID :exec
DELETE FROM events
WHERE user_id = ?;
-- name: GetAllEventsByUserID :many
SELECT *
FROM events
WHERE user_id = ?;
-- name: GetEventByTitle :one
SELECT *
FROM events
WHERE title = ?;
-- name: GetEventByID :one
SELECT *
FROM events
WHERE event_id = ?;

