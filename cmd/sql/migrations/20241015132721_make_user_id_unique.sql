-- +goose Up
-- +goose StatementBegin
ALTER TABLE `users`
ADD UNIQUE (`user_id`);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE `users`
ADD UNIQUE `user_id`;
-- +goose StatementEnd