-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
ADD CONSTRAINT UNIQUE (event_id);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE attendees
ADD event_id INT NOT NULL,
ADD CONSTRAINT fk_attendees_events
FOREIGN KEY (event_id) REFERENCES events (event_id)
ON UPDATE CASCADE
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE attendees
DROP FOREIGN KEY fk_attendees_events,
DROP COLUMN event_id;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE events
DROP INDEX event_id;
-- +goose StatementEnd
