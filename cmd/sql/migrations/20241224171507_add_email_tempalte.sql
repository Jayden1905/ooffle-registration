-- +goose Up
CREATE TABLE IF NOT EXISTS `email_template` (
    `id` int NOT NULL AUTO_INCREMENT,
    `event_id` int NOT NULL,
    `header_image` text NOT NULL,
    `content` text NOT NULL,
    `footer_image` text NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `fk_email_template_events` (`event_id`),
    CONSTRAINT `fk_email_template_events` FOREIGN KEY (`event_id`) REFERENCES `events` (`event_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
-- +goose Down
DROP TABLE `email_template`;