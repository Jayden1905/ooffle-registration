-- +goose Up
-- +goose StatementBegin
CREATE TABLE `attendees` (
  `id` int NOT NULL AUTO_INCREMENT,
  `first_name` varchar(50) NOT NULL,
  `last_name` varchar(50) NOT NULL,
  `email` varchar(255) NOT NULL,
  `qr_code` text DEFAULT NULL,
  `company_name` varchar(50) DEFAULT NULL,
  `title` varchar(50) DEFAULT NULL,
  `table_no` int DEFAULT NULL,
  `role` varchar(50) DEFAULT NULL,
  `attendance` enum('Yes', 'No') DEFAULT 'No',
  PRIMARY KEY (`id`),
  UNIQUE KEY `Email_UNIQUE` (`email`),
  KEY `idx_lastName` (`last_name`),
  KEY `idx_company` (`company_name`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `attendees`;
-- +goose StatementEnd
