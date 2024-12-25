-- +goose Up
-- +goose StatementBegin
ALTER TABLE `email_template`
ADD CONSTRAINT `email_template_event_id_unique` UNIQUE (`event_id`);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE `email_template` DROP INDEX `email_template_event_id_unique`;
-- +goose StatementEnd