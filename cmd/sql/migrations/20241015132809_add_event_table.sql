-- +goose Up
CREATE TABLE IF NOT EXISTS `events` (
    `event_id` int NOT NULL AUTO_INCREMENT,
    `title` varchar(50) NOT NULL,
    `description` TEXT NOT NULL,
    `start_date` datetime NOT NULL,
    `end_date` datetime NOT NULL,
    `location` varchar(255) NOT NULL,
    `user_id` int NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`event_id`, `user_id`),
    KEY `fk_events_users1_idx` (`user_id`),
    CONSTRAINT `fk_events_users1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
-- +goose Down
DROP TABLE IF EXISTS `events`;
