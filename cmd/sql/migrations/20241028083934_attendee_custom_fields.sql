-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `attendees_custom_fields` (
  `id` int NOT NULL AUTO_INCREMENT,
  `attendee_id` int NOT NULL,
  `field_name` varchar(100) DEFAULT NULL,
  `field_value` text,
  `field_type` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_attendeesCustomFields_attendees` (`attendee_id`),
  CONSTRAINT `attendees_custom_fields_ibfk_1` FOREIGN KEY (`attendee_id`) REFERENCES `attendees` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `attendees_custom_fields`;
-- +goose StatementEnd
