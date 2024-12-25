-- +goose Up
-- +goose StatementBegin
ALTER TABLE attendees DROP INDEX Email_UNIQUE;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE attendees DROP INDEX Email_UNIQUE;
-- +goose StatementEnd