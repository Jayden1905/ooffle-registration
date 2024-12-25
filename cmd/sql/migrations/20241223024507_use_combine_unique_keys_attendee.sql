-- +goose Up
-- +goose StatementBegin
ALTER TABLE attendees
ADD UNIQUE KEY Email_Event_UNIQUE (email, event_id);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE attendees DROP INDEX Email_Event_UNIQUE;
-- +goose StatementEnd