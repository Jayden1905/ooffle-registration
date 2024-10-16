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
-- name: GetAllEvents :many
SELECT *
FROM events;
-- name: GetEventByTitle :one
SELECT *
FROM events
WHERE title = ?;
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