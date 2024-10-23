-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN verify BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN verify;
-- +goose StatementEnd