-- +goose Up
-- +goose StatementBegin
ALTER TABLE `email_template`
ADD COLUMN `subject` text,
    ADD COLUMN `message` text,
    ADD COLUMN `bg_color` text;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE `email_template` DROP COLUMN `subject`,
    DROP COLUMN `message`,
    DROP COLUMN `bg_color`;
-- +goose StatementEnd